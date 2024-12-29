package apiserver

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"iam/internal/apiserver/config"
	"iam/internal/apiserver/store"
	"iam/internal/apiserver/store/mysql"
	"iam/internal/pkg/options"
	genericserver "iam/internal/pkg/server"
	"iam/pkg/cache"
	"iam/pkg/shutdown"
	"iam/pkg/shutdown/shutdownmanagers/posixsignal"
	"log"
)

type apiServer struct {
	gs    *shutdown.GracefulShutdown
	http  *options.InsecureServingOptions
	https *options.SecureServingOptions
	grpc  *options.GrpcOptions // 来自config

	mysql         *options.MysqlOptions // 初始化mysql所需
	redis         *options.RedisOptions
	jwt           *options.JwtOptions
	GrpcServer    *genericserver.GrpcAPIServer
	GenericServer *genericserver.GenericAPIServer
	serverRun     *options.ServerRunOptions // mode healthz middleware  apply 进行构建到pkg.config中
	feature       *options.FeatureOptions   // pprof metrics apply 进行构建到pkg.config中
	log           *options.LogOption
}

// 对apiServer进行相关准备工作
type perparesApiServer struct {
	*apiServer
}

// 将 options 中配置参数加载到 generics_server 中
func (server *apiServer) initApplyServer() (*genericserver.Config, error) {

	var (
		genericConfig = genericserver.NewConfig()
	)

	if err := server.http.ApplyTo(genericConfig); err != nil {
		return nil, err
	}

	if err := server.https.ApplyTo(genericConfig); err != nil {
		return nil, err
	}
	if err := server.jwt.ApplyTo(genericConfig); err != nil {
		return nil, err
	}

	if err := server.serverRun.ApplyTo(genericConfig); err != nil {
		return nil, err
	}

	if err := server.feature.ApplyTo(genericConfig); err != nil {
		return nil, err
	}

	if err := server.jwt.ApplyTo(genericConfig); err != nil {
		return nil, err
	}

	// 未应用 http

	return genericConfig, nil
}

// 构建http/grpc服务
func newApiServer(cfg *config.Config) (*apiServer, error) {

	// 创建 shutdown 实例
	gs := shutdown.New()
	// 添加信号
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager()) // os.Interrupt , syscall.SIGTERM 的信号

	var (
		err       error
		serverCfg *genericserver.Config
	)

	server := &apiServer{
		gs:        gs,
		http:      cfg.Http,
		https:     cfg.Https,
		grpc:      cfg.Grpc,
		mysql:     cfg.Mysql,
		redis:     cfg.Redis,
		jwt:       cfg.Jwt,
		serverRun: cfg.ServerRun,
		feature:   cfg.Feature,
	}

	// 应用 generic server 所需配置参数
	if serverCfg, err = server.initApplyServer(); err != nil {
		return nil, err
	}

	// serverCfg -》 这里调用的New方法，这种设计方法可以确保我们创建的服务实例一定是基于 complete 之后的配置（completed） ，
	// Complete进行构建实例返回的为补全后的配置信息 --> 配置中多了 http.Server 等
	if server.GenericServer, err = serverCfg.Complete().New(); err != nil {
		return nil, err
	}
	cfg.Grpc.ServerCert = cfg.Https.ServerCert

	// --- 构建 grpc 服务配置 ----
	grpcCfg := genericserver.NewGrpcConfig()
	// 应用 grpc 配置
	if err = cfg.Grpc.Apply(grpcCfg); err != nil {
		return nil, err
	}
	// 构建 grpc 配置
	server.GrpcServer = grpcCfg.NewCompletedGrpc().New()

	return server, nil
}

// PrepareServer 准备服务
func (server *apiServer) PrepareServer() *perparesApiServer {

	// 初始化 redis mysql
	// server.initRedisStore() --> 暂时不进行测试
	// server.initMysqlStore()

	// es
	// mq

	// 添加出现对应信号时候 进行执行相关回调函数 （关闭连接等）
	server.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		mysqlStore, _ := mysql.GetStoreFactory(nil)
		// 关闭 mysql
		if mysqlStore != nil {
			_ = mysqlStore.Close()
		}
		server.GrpcServer.Close()
		server.GenericServer.Close()
		return nil
	}))

	// 构建路由
	initRouter(server.GenericServer.Engine)

	// todo 注册 pb 服务
	// pb.RegisterCacheServer(server.GrpcServer.Server, cacheIns)

	return &perparesApiServer{server}
}

// Run 实际运行
func (preSvc *perparesApiServer) Run() error {

	// 启动 Shutdown 服务
	if err := preSvc.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	// grpc 运行
	go preSvc.GrpcServer.Run()

	// http/https 运行
	return preSvc.GenericServer.Run()
}

// initShutdown 需要 --》 shutdown mysql redis http grpc

func (server *apiServer) initMysqlStore() {
	logrus.Debug("初始化mysql")
	factory, err := mysql.GetStoreFactory(server.mysql)
	if err != nil {
		panic(fmt.Errorf("init mysql err:%v", err)) // -->
	}

	store.SetFactory(factory)

}

// 初始化redis
func (server *apiServer) initRedisStore() {

	ctx, cancel := context.WithCancel(context.Background())
	server.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		cancel()

		return nil
	}))

	cfg := &cache.Config{
		Host:                  server.redis.Host,
		Port:                  server.redis.Port,
		Addrs:                 server.redis.Addrs,
		MasterName:            server.redis.MasterName,
		Username:              server.redis.Username,
		Password:              server.redis.Password,
		Database:              server.redis.Database,
		MaxIdle:               server.redis.MaxIdle,
		MaxActive:             server.redis.MaxActive,
		Timeout:               server.redis.Timeout,
		EnableCluster:         server.redis.EnableCluster,
		UseSSL:                server.redis.UseSSL,
		SSLInsecureSkipVerify: server.redis.SSLInsecureSkipVerify,
	}

	// 初始化 redis 并尝试重连
	go cache.ConnectToRedis(ctx, cfg)
}
