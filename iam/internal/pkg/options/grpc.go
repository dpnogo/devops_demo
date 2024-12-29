package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"iam/internal/pkg/server"
	"net"
	"strconv"
)

type GrpcOptions struct {
	BindAddress string   `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int      `json:"bind-port" mapstructure:"bind-port"`
	MaxBodySize int      `json:"max-body-size" mapstructure:"max-body-size"`
	ServerCert  CertConf `json:"tls"` // 使用为https的tls
}

func NewGrpcOptions() *GrpcOptions {
	return &GrpcOptions{
		BindAddress: "127.0.0.1",
		BindPort:    8081,
		MaxBodySize: 4 * 1024 * 1024,
	}
}

func (s *GrpcOptions) Validate() []error {
	var errors []error

	if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--insecure-port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port",
				s.BindPort,
			),
		)
	}
	return errors
}

func (s *GrpcOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.BindAddress, "grpc.bind-address", s.BindAddress, "Grpc server bind address.")
	fs.StringVar(&s.BindAddress, "grpc.bind-port", s.BindAddress, "Grpc server bind port.")
	fs.IntVar(&s.MaxBodySize, "grpc.max-body-size", s.MaxBodySize, "Grpc server request body max size.")
}

// Apply 应用grpc server config
func (s *GrpcOptions) Apply(grpcCfg *server.GrpcConfig) error {
	grpcCfg.MaxMsgSize = s.MaxBodySize
	grpcCfg.Addr = net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))

	grpcCfg.CertFile = s.ServerCert.CertKey.CertFile
	grpcCfg.KeyFile = s.ServerCert.CertKey.KeyFile

	return nil
}
