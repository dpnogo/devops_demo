package options

import (
	"iam/internal/authzserver/analytics"
	genericoptions "iam/internal/pkg/options"
	"iam/internal/pkg/server"
	pkg "iam/pkg/app/cli"
)

// http
// grpc
// redis
// pg
// 消费日志

// 从 pkg 中 选择 authzserver 中需要 options , 应用配置

// Options 实现 apps 对应的CliOptions接口
type Options struct {
	SecureServing    *genericoptions.SecureServingOptions   `json:"secure"`                                  // https 配置
	InsecureServing  *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`        // http 配置
	Feature          *genericoptions.FeatureOptions         `json:"feature" mapstructure:"feature"`          // 运行设置 功能 是否开启相关分析功能
	ServerRun        *genericoptions.ServerRunOptions       `json:"server" mapstructure:"server"`            // 运行时配置
	Jwt              *genericoptions.JwtOptions             `json:"jwt" mapstructure:"jwt"`                  // jwt 配置
	AnalyticsOptions *analytics.AnalyticsOptions            `json:"analytics"      mapstructure:"analytics"` // 授权日志写到redis中配置
	RedisOptions     *genericoptions.RedisOptions           `json:"redis"          mapstructure:"redis"`
	// Log                     *log.Options                           `json:"log"            mapstructure:"log"`
	RPCServer string `json:"rpcserver"      mapstructure:"rpcserver"` // authz只需要调用api-server所以仅仅只需要
	ClientCA  string `json:"client-ca-file" mapstructure:"client-ca-file"`
}

// NewOptions 创建option的默认配置
func NewOptions() *Options {
	return &Options{
		RPCServer:        "127.0.0.1:8081",
		ClientCA:         "",
		SecureServing:    genericoptions.NewSecureServing(),
		InsecureServing:  genericoptions.NewInsecureServingOptions(),
		Feature:          genericoptions.NewFeatureOptions(),
		ServerRun:        genericoptions.NewServerRunOptions(),
		Jwt:              genericoptions.NewJwtOptions(), // 注:按正常来说这里应该是走
		AnalyticsOptions: analytics.NewAnalyticsOptions(),
		RedisOptions:     genericoptions.NewRedisOptions(),
	}
}

// ApplyTo 应用其他配置 , 暂时可有可无（这里）
func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

// Flags 创建flag name 并初始化获取所有的flag值
func (o *Options) Flags() (fss pkg.NamedFlagSets) {

	o.SecureServing.AddFlags(fss.FlagSet("secure"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure"))
	o.Feature.AddFlags(fss.FlagSet("feature"))
	o.ServerRun.AddFlags(fss.FlagSet("server"))
	o.Jwt.AddFlags(fss.FlagSet("jwt"))
	o.AnalyticsOptions.AddFlags(fss.FlagSet("analytics"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))

	// 其他非构建的杂项
	fs := fss.FlagSet("misc")
	fs.StringVar(&o.RPCServer, "rpcserver", o.RPCServer, "dial apiserver grpc addr")
	fs.StringVar(&o.RPCServer, "client-ca-file", o.RPCServer, "dial apiserver grpc addr")

	return fss
}

func (o *Options) Validate() []error {

	var errs = make([]error, 0)
	errs = append(errs, o.SecureServing.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	errs = append(errs, o.Feature.Validate()...)
	errs = append(errs, o.ServerRun.Validate()...)
	errs = append(errs, o.Jwt.Validate()...)
	errs = append(errs, o.AnalyticsOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)

	return errs
}

// viper 读取配置的优先级， 高优先级配置会进行覆盖低优先级的配置
// 从高到低分别为
// 通过 viper.Set 函数显示设置的配置
// 命令行参数
// 环境变量
// 配置文件
// Key/Value 存储
// 默认值

// 验证得到参数的配置 -->
