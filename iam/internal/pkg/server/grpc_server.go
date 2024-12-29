package server

import (
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcAPIServer struct {
	*grpc.Server
	address string
}

// Run 运行 ,若存在stop信号时
func (grpcSvr *GrpcAPIServer) Run() {

	log.Printf("start grpc run:%s\n", grpcSvr.address)

	listen, err := net.Listen("tcp", grpcSvr.address)
	if err != nil {
		log.Fatalf("failed to start grpc server: %s", err.Error())
	}

	go func() {
		if err := grpcSvr.Serve(listen); err != nil {
			log.Fatalf("failed to start grpc server: %s", err.Error())
		}
	}()

	log.Println("start grpc complete")

}

func (grpcSvr *GrpcAPIServer) Close() {
	grpcSvr.GracefulStop()
	log.Printf("GRPC server on %s stopped", grpcSvr.address)
}
