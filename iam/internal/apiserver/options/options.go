package options

import (
	"encoding/json"
	"iam/internal/pkg/options"
	pkg "iam/pkg/app/cli"
	"iam/pkg/util/idutil"
)

// Options 实现 CliOptions 接口以提供 apps 中使用
//type CliOptions interface {
//	Flags() pkg.NamedFlagSets // 获取所有的flag值
//	Validate() []error        // 对相关的配置进行校验 同时验证所有的命令行参数/配置文件是否有效
//}

// Options .keep-server 所需配置
type Options struct {
	Http  *options.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	Https *options.SecureServingOptions   `json:"secure" mapstructure:"secure"`
	Grpc  *options.GrpcOptions            `json:"grpc" mapstructure:"grpc"`
	Mysql *options.MysqlOptions           `json:"mysql" mapstructure:"mysql"`
	Redis *options.RedisOptions           `json:"redis" mapstructure:"redis"`
	Jwt   *options.JwtOptions             `json:"jwt" mapstructure:"jwt"`
	// 服务相关配置
	ServerRun *options.ServerRunOptions `json:"server" mapstructure:"server"`   // mode healthz middleware  apply 进行构建到pkg.config中
	Feature   *options.FeatureOptions   `json:"feature" mapstructure:"feature"` // pprof metrics apply 进行构建到pkg.config中
	Log       *options.LogOption        `json:"log" mapstructure:"log"`         // 日志
}

// NewOptions 设置默认配置
func NewOptions() *Options {
	newOp := &Options{
		Http:      options.NewInsecureServingOptions(),
		Https:     options.NewSecureServing(),
		Grpc:      options.NewGrpcOptions(),
		Jwt:       options.NewJwtOptions(),
		Mysql:     options.NewMysqlOptions(),
		Redis:     options.NewRedisOptions(), // 更新或者其他操作策略，进行更改
		ServerRun: options.NewServerRunOptions(),
		Feature:   options.NewFeatureOptions(),
		Log:       options.NewLogOption(),
	}
	return newOp
}

// Flags 添加对应flag分类
func (ops *Options) Flags() (fss pkg.NamedFlagSets) {

	ops.Http.AddFlags(fss.FlagSet("http"))
	ops.Https.AddFlags(fss.FlagSet("https"))
	ops.Grpc.AddFlags(fss.FlagSet("grpc"))
	ops.Mysql.AddFlags(fss.FlagSet("mysql"))
	ops.Redis.AddFlags(fss.FlagSet("redis"))
	ops.Jwt.AddFlags(fss.FlagSet("jwt"))
	ops.ServerRun.AddFlags(fss.FlagSet("server"))
	ops.Feature.AddFlags(fss.FlagSet("feature"))
	ops.Log.AddFlags(fss.FlagSet("logger"))

	return fss
}

// Validate 验证参数是否存在问题
func (ops *Options) Validate() []error {

	var errs = make([]error, 0)
	errs = append(errs, ops.Http.Validate()...)
	errs = append(errs, ops.Https.Validate()...)
	errs = append(errs, ops.Grpc.Validate()...)
	errs = append(errs, ops.Mysql.Validate()...)
	errs = append(errs, ops.Redis.Validate()...)
	errs = append(errs, ops.Jwt.Validate()...)
	errs = append(errs, ops.ServerRun.Validate()...)
	errs = append(errs, ops.Feature.Validate()...)
	errs = append(errs, ops.Log.Validate()...)

	return errs
}

// CompleteableOptions abstracts options which can be completed.
//type CompleteableOptions interface {
//	Complete() error
//}

// Complete 实现 app 的 CompleteableOptions 接口 : 设置默认选项。
func (ops *Options) Complete() error {
	if ops.Jwt.Key == "" {
		ops.Jwt.Key = idutil.NewSecretKey() // 若不指定，则随机key
	}

	return ops.Https.Complete() // 验证 https tls的，并根据配置/参数进行选择
}

// PrintableOptions abstracts options which can be printed.
//type PrintableOptions interface {
//	String() string
//}

// 实现 app 的 PrintableOptions 接口 : 抽象可打印的接口
func (ops *Options) String() string {
	data, _ := json.Marshal(ops)
	return string(data)
}
