package options

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"iam/internal/pkg/server"
)

// ServerRunOptions 服务允许途中需要的配置
type ServerRunOptions struct {
	Mode        string   `json:"mode"        mapstructure:"mode"`
	Healthz     bool     `json:"healthz"     mapstructure:"healthz"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

func NewServerRunOptions() *ServerRunOptions {
	defaultCfg := server.NewConfig()

	return &ServerRunOptions{
		Mode:        defaultCfg.Mode,
		Healthz:     defaultCfg.Healthz,
		Middlewares: defaultCfg.Middlewares,
	}
}

// Validate 验证服务运行时配置参数
func (option *ServerRunOptions) Validate() []error {

	var err = []error{}

	switch option.Mode {
	case gin.ReleaseMode, gin.DebugMode, gin.TestMode:

	default:
		err = append(err, fmt.Errorf("gin mode must dev,test,release type, mode:%s", option.Mode))
	}
	return err
}

func (option *ServerRunOptions) ApplyTo(config *server.Config) error {
	config.Mode = option.Mode
	config.Healthz = option.Healthz
	config.Middlewares = option.Middlewares
	return nil
}

func (option *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&option.Mode, "mode", option.Mode, "gin run mode type")
	fs.BoolVar(&option.Healthz, "healthz", option.Healthz, "server add or not health check")
	fs.StringSliceVar(&option.Middlewares, "middlewares", option.Middlewares, "server use some middleware")
}
