package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

// GrpcConfig grpc 从internal的服务中选择需要的配置 -->
type GrpcConfig struct {
	Addr       string // 地址
	MaxMsgSize int    // 最大msg

	// tls相关 防止和option引用圈
	CertFile string
	KeyFile  string
}

// NewGrpcConfig 创建默认配置
func NewGrpcConfig() *GrpcConfig {
	return &GrpcConfig{
		Addr:       "127.0.0.1:8081",
		MaxMsgSize: 1024 * 1024 * 10,
	}
}

type CompletedGrpc struct {
	*GrpcConfig
}

func (c *GrpcConfig) NewCompletedGrpc() *CompletedGrpc {
	return &CompletedGrpc{c}
}

// New 构建run
func (grpcSvr *CompletedGrpc) New() *GrpcAPIServer {

	creds, err := credentials.NewServerTLSFromFile(grpcSvr.CertFile, grpcSvr.KeyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err.Error())
	}

	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(grpcSvr.MaxMsgSize), grpc.Creds(creds)}

	grpcServer := grpc.NewServer(opts...)

	return &GrpcAPIServer{
		grpcServer,
		grpcSvr.Addr,
	}
}

/*
creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err.Error())
	}
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)

   client
	storeIns, _ := mysql.GetMySQLFactoryOr(c.mysqlOptions)
	store.SetClient(storeIns)
   redis
	cacheIns, err := cachev1.GetCacheInsOr(storeIns)
	if err != nil {
		log.Fatalf("Failed to get cache instance: %s", err.Error())
	}

	pb.RegisterCacheServer(grpcServer, cacheIns)

	reflection.Register(grpcServer)

	return &grpcAPIServer{grpcServer, c.Addr}, nil
*/
