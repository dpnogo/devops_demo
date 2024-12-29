package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"iam/internal/pkg/server"
	"net"
	"strconv"
)

// InsecureServingOptions 不安全服务选项
type InsecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port" mapstructure:"bind-port"`
}

// NewInsecureServingOptions 初始化服务
func NewInsecureServingOptions() *InsecureServingOptions {
	return &InsecureServingOptions{
		BindAddress: "127.0.0.1",
		BindPort:    8080,
	}
}

func (svc *InsecureServingOptions) Validate() []error {

	if svc.BindPort < 0 || svc.BindPort > 65535 {
		var eList = make([]error, 0)
		eList = append(eList, fmt.Errorf("port must 0 <= port <= 65535 , port:%d", svc.BindPort))
		return eList
	}

	return []error{}
}

// ApplyTo 构建address
func (svc *InsecureServingOptions) ApplyTo(config *server.Config) error {
	// 将 svc 应用到 server的config中供svc使用
	config.InsecureServing = &server.InsecureServingInfo{
		Address: net.JoinHostPort(svc.BindAddress, strconv.Itoa(svc.BindPort)),
	}
	return nil
}

func (svc *InsecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&svc.BindAddress, "address", "", "insecure_svc_options address flag") // flag 对应的值
	fs.IntVarP(&svc.BindPort, "port", "p", 0, "insecure_svc_options port flag, port must >= 0 and <= 65535")
}
