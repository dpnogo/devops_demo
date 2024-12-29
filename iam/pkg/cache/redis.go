package cache

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/go-redis/redis"
	"log"
	"sync/atomic"
	"time"
)

var (
	singlePool      atomic.Value // 存储 redis client
	singleCachePool atomic.Value // cache -->  redis client 连接池
	redisUp         atomic.Value // 连接redis，此时是否可用
)

// ErrRedisIsDown Redis is either down or ws not configured
var ErrRedisIsDown = errors.New("storage: Redis is either down or ws not configured")

// Config defines options for redis cluster.
type Config struct {
	Host                  string
	Port                  int
	Addrs                 []string
	MasterName            string
	Username              string
	Password              string
	Database              int
	MaxIdle               int
	MaxActive             int
	Timeout               int
	EnableCluster         bool
	UseSSL                bool
	SSLInsecureSkipVerify bool
}

// RedisOpts redis 连接所需配置
type RedisOpts redis.UniversalOptions

// NewRedisClusterPool 创建 redis 缓存池
func NewRedisClusterPool(isCache bool, config *Config) redis.UniversalClient {

	// 创建连接大小
	poolSize := 500
	if config.MaxActive > 0 {
		poolSize = config.MaxActive
	}

	// 连接超时
	timeout := 5 * time.Second
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	// 通过 tls 方式连接redis
	var tlsConfig *tls.Config
	if config.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}
	var client redis.UniversalClient

	opts := RedisOpts{
		Addrs:        config.Addrs,
		MasterName:   config.MasterName,
		Password:     config.Password,
		DB:           config.Database,
		DialTimeout:  timeout, // 连接超时时间
		ReadTimeout:  timeout, // 读取超时时间
		WriteTimeout: timeout,
		IdleTimeout:  240 * timeout,
		PoolSize:     poolSize,
		TLSConfig:    tlsConfig,
	}

	// 根据配置选择不同的初始化方式 , 使用主从服务方式进行连接
	if config.MasterName != "" {
		log.Println("dial redis by type master-slave")
		client = redis.NewFailoverClient(opts.failoverOpt())
	} else if config.EnableCluster { // 说明为 cluster 模式
		client = redis.NewClusterClient(opts.cluster())
	} else {
		client = redis.NewClient(opts.normal())
	}
	return client
}

// 构建主从模式所需配置
func (opts *RedisOpts) failoverOpt() *redis.FailoverOptions {

	if len(opts.Addrs) == 0 {
		opts.Addrs = []string{"127.0.0.1:26379"}
	}

	return &redis.FailoverOptions{
		MasterName:    opts.MasterName,
		SentinelAddrs: opts.Addrs, // Sentinel 地址 哨兵地址

		OnConnect: opts.OnConnect, //

		DB:       opts.DB,
		Password: opts.Password,

		MaxRetries:      opts.MaxRetries,
		MinRetryBackoff: opts.MinRetryBackoff,
		MaxRetryBackoff: opts.MaxRetryBackoff,

		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,

		PoolSize:           opts.PoolSize,
		MinIdleConns:       opts.MinIdleConns,
		MaxConnAge:         opts.MaxConnAge,
		PoolTimeout:        opts.PoolTimeout,
		IdleTimeout:        opts.IdleTimeout,
		IdleCheckFrequency: opts.IdleCheckFrequency,

		TLSConfig: opts.TLSConfig,
	}
}

// 构建集群所需配置
func (opts *RedisOpts) cluster() *redis.ClusterOptions {
	if len(opts.Addrs) == 0 {
		opts.Addrs = []string{"127.0.0.1:26379"}
	}

	return &redis.ClusterOptions{
		Addrs:     opts.Addrs,
		OnConnect: opts.OnConnect,

		Password: opts.Password,

		MaxRedirects:   opts.MaxRedirects,
		ReadOnly:       opts.ReadOnly,
		RouteByLatency: opts.RouteByLatency,
		RouteRandomly:  opts.RouteRandomly,

		MaxRetries:      opts.MaxRetries,
		MinRetryBackoff: opts.MinRetryBackoff,
		MaxRetryBackoff: opts.MaxRetryBackoff,

		DialTimeout:        opts.DialTimeout,
		ReadTimeout:        opts.ReadTimeout,
		WriteTimeout:       opts.WriteTimeout,
		PoolSize:           opts.PoolSize,
		MinIdleConns:       opts.MinIdleConns,
		MaxConnAge:         opts.MaxConnAge,
		PoolTimeout:        opts.PoolTimeout,
		IdleTimeout:        opts.IdleTimeout,
		IdleCheckFrequency: opts.IdleCheckFrequency,

		TLSConfig: opts.TLSConfig,
	}
}

// 构建单机所需配置
func (opts *RedisOpts) normal() *redis.Options {

	addr := "127.0.0.1:6379"
	if len(opts.Addrs) > 0 {
		addr = opts.Addrs[0]
	}

	return &redis.Options{
		Addr:      addr,
		OnConnect: opts.OnConnect,

		DB:       opts.DB,
		Password: opts.Password,

		MaxRetries:      opts.MaxRetries,
		MinRetryBackoff: opts.MinRetryBackoff,
		MaxRetryBackoff: opts.MaxRetryBackoff,

		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,

		PoolSize:           opts.PoolSize,
		MinIdleConns:       opts.MinIdleConns,
		MaxConnAge:         opts.MaxConnAge,
		PoolTimeout:        opts.PoolTimeout,
		IdleTimeout:        opts.IdleTimeout,
		IdleCheckFrequency: opts.IdleCheckFrequency,

		TLSConfig: opts.TLSConfig,
	}
}

// 此时redis是否可用
func (r *RedisCluster) up() error {
	if !Connected() {
		return ErrRedisIsDown
	}
	return nil
}

func Connected() bool {
	if v := redisUp.Load(); v != nil {
		return v.(bool)
	}
	return false
}

// RedisCluster 是一个使用redis数据库的存储管理器
type RedisCluster struct {
	KeyPrefix string
	HashKeys  bool
	IsCache   bool
}

// ConnectToRedis 若redis断开，则定时尝试连接redis
func ConnectToRedis(ctx context.Context, config *Config) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	c := []RedisCluster{
		{}, {IsCache: true},
	}
	var ok bool
	for _, v := range c {
		// 查看缓存池是否存在client，若无新进行创建
		if !connectSingleton(v.IsCache, config) {
			break
		}
		// 尝试打开一个,若无可用连接，或者连接不可用
		if !clusterConnectionIsOpen(v) {
			redisUp.Store(false) // 设置redis状态不可用

			break
		}
		ok = true
	}
	redisUp.Store(ok) // 设置redis连接状态

again:
	// 定时重连, 即重新创建 client,并添加到
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if !shouldConnect() { // 若禁用redis，则直接 continue --> 暂时无用
				continue
			}
			for _, v := range c {
				// 查看缓存池是否存在client，若无新进行创建
				if !connectSingleton(v.IsCache, config) {
					redisUp.Store(false)
					goto again
				}
				// 尝试打开一个,若无可用连接，或者连接不可用
				if !clusterConnectionIsOpen(v) {
					redisUp.Store(false)

					goto again
				}
			}
			redisUp.Store(true) // redis可用
		}
	}
}

var disableRedis atomic.Value // 禁用redis

func shouldConnect() bool {
	if v := disableRedis.Load(); v != nil {
		return !v.(bool)
	}

	return true
}

// 查询 singlePool || singleCachePool 是否存在
func singleton(cache bool) redis.UniversalClient {

	if cache {
		v := singleCachePool.Load()
		if v != nil {
			return v.(redis.UniversalClient)
		}
		return nil
	}

	v := singlePool.Load()
	if v != nil {
		return v.(redis.UniversalClient)
	}

	return nil
}

// 若无，则进行初始化client到 singlePool || singleCachePool 中
func connectSingleton(cache bool, cfg *Config) bool {

	// 说明此时 singlePool || singleCachePool 无值，则进行创建
	if singleton(cache) == nil {
		if cache {
			singleCachePool.Store(NewRedisClusterPool(cache, cfg))
			return true
		}
		singlePool.Store(NewRedisClusterPool(cache, cfg))
		return true
	}
	return true
}

// clusterConnectionIsOpen 是否能够正常连接
func clusterConnectionIsOpen(cluster RedisCluster) bool {
	c := singleton(cluster.IsCache) // 获取一个client
	// 并进行测试是否能够正常连通

	key := "cluster-key"
	val := "cluster-val"

	if err := c.Set(key, val, time.Second).Err(); err != nil {
		return false
	}

	if err := c.Get(key).Err(); err != nil {
		return false
	}

	return true
}

// ---------- 所使用的API ----------

// 得到一个client
func (r *RedisCluster) singleton() redis.UniversalClient {
	return singleton(r.IsCache)
}

// Publish 往 redis push
func (r *RedisCluster) Publish(psChannel, message string) error {

	if err := r.up(); err != nil {
		return err
	}

	err := r.singleton().Publish(psChannel, message).Err()
	if err != nil {
		return err
	}

	return nil
}
