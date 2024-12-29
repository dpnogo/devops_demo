package cache

import (
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/ory/ladon"
	"github.com/pkg/errors"
	"iam/internal/authzserver/store"
	pb "iam/pkg/proto/apiserver/v1"
	"sync"
)

var (
	ErrSecretNotFound = fmt.Errorf("secret not found")
	ErrPolicyNotFound = fmt.Errorf("secret not found")
)

// cache 实现 load 接口

type Cache struct {
	lock     *sync.RWMutex
	cli      store.Factory
	secrets  *ristretto.Cache
	policies *ristretto.Cache
}

var (
	cacheIns *Cache // 缓冲
	syncOnce sync.Once
)

// GetCacheInsOr 初始化缓存库
func GetCacheInsOr(cli store.Factory) (*Cache, error) {

	var err error

	if cli != nil {
		var (
			secretCache *ristretto.Cache
			policyCache *ristretto.Cache
		)
		syncOnce.Do(func() {
			c := &ristretto.Config{
				NumCounters: 1e7,     // number of keys to track frequency of (10M).
				MaxCost:     1 << 30, // maximum cost of cache (1GB).
				BufferItems: 64,      // number of keys per Get buffer.
				Cost:        nil,
			}
			secretCache, err = ristretto.NewCache(c)
			if err != nil {
				return
			}
			policyCache, err = ristretto.NewCache(c)
			if err != nil {
				return
			}

			cacheIns = &Cache{
				secrets:  secretCache,
				policies: policyCache,
				lock:     new(sync.RWMutex),
				cli:      cli,
			}

		})
	}

	return cacheIns, err
}

// GetSecret 从缓存中获取到对应的密钥信息
func (c *Cache) GetSecret(key string) (*pb.SecretInfo, error) {

	c.lock.RLock()
	defer c.lock.RUnlock()

	value, ok := c.secrets.Get(key)
	if !ok {
		return nil, ErrSecretNotFound
	}

	return value.(*pb.SecretInfo), nil
}

// GetPolicy 从缓存中获取到对应的策略信息
func (c *Cache) GetPolicy(key string) ([]*ladon.DefaultPolicy, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	value, ok := c.policies.Get(key)
	if !ok {
		return nil, ErrPolicyNotFound
	}

	return value.([]*ladon.DefaultPolicy), nil
}

// Reload 重新加载全部的密钥和策略， 重新加载时候（开始运行 + 定时）
func (c *Cache) Reload() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// reload secrets
	secrets, err := c.cli.Secrets().List()
	if err != nil {
		return errors.Wrap(err, "list secrets failed")
	}

	c.secrets.Clear()
	for key, val := range secrets {
		c.secrets.Set(key, val, 1)
	}

	// reload policies
	policies, err := c.cli.Policies().List()
	if err != nil {
		return errors.Wrap(err, "list policies failed")
	}

	c.policies.Clear()
	for key, val := range policies {
		c.policies.Set(key, val, 1)
	}

	return nil
}
