package authzserver

import (
	"context"
	"iam/internal/authzserver/analytics"
	"iam/internal/authzserver/config"
	genericoptions "iam/internal/pkg/options"
	genericserver "iam/internal/pkg/server"
	"iam/pkg/shutdown"
	"iam/pkg/shutdown/shutdownmanagers/posixsignal"
	"log"
)

// 根据 config 进行构建authz服务

/*
 config -> prepared Run --> 实际运行
*/

type authzServer struct {
	gs               *shutdown.GracefulShutdown // 启动/结束 时候需要回调函数
	rpcServer        string
	clientCA         string
	redisOptions     *genericoptions.RedisOptions
	genericAPIServer *genericserver.GenericAPIServer // 部分功能抽离到 pkg.server中，构建http服务
	analyticsOptions *analytics.AnalyticsOptions
	redisCancelFunc  context.CancelFunc // redis 回调函数
}

type preparedServer struct {
	*authzServer
}

// createAuthzServer 通过 authorization 的 config 初始化 authzServer
func createAuthzServer(cfg *config.Config) (*authzServer, error) {

	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager()) // todo ---->

	authSvc := &authzServer{
		gs:           shutdown.New(),
		rpcServer:    cfg.RPCServer,
		clientCA:     cfg.ClientCA,
		redisOptions: cfg.RedisOptions,
	}

	// 加载
	// cfg.Option.ApplyTo(svcCfg)
	svcCfg, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}
	// server config -> server GenericAPIServer Complete() 处理或补全配置
	authSvc.genericAPIServer, _ = svcCfg.Complete().New()

	return authSvc, nil
}

// 通过应用配置 构建 HTTP/GRPC 服务配置
func buildGenericConfig(cfg *config.Config) (genericConfig *genericserver.Config, lastErr error) {
	genericConfig = genericserver.NewConfig()
	if lastErr = cfg.ServerRun.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.Feature.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

func (svc *authzServer) PreparedServer() *preparedServer {

	// 初始化redis + 初始化 将 .keep-server的策略和密钥 初始化到内存中

	// 初始化 router
	initRouter(svc.genericAPIServer.Engine)

	return &preparedServer{svc}
}

// Run 前置条件准备完成后，实际运行
func (preSvc *preparedServer) Run() error {

	preSvc.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		preSvc.genericAPIServer.Close()
		if preSvc.analyticsOptions.Enable {
			analytics.GetAnalytics().Stop()
		}
		preSvc.redisCancelFunc()

		return nil
	}))

	// start shutdown managers
	if err := preSvc.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	// 运行 HTTP 服务
	return preSvc.genericAPIServer.Run()
}

// init 密钥和策略 -->
