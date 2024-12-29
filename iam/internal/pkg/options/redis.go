package options

import (
	"github.com/spf13/pflag"
)

type RedisOptions struct {
	Host     string `json:"host" mapstructure:"host" description:"Redis service host address"`
	Port     int    `json:"port"  mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Database int    `json:"database" mapstructure:"database"`

	Addrs                 []string `json:"addrs"                    mapstructure:"addrs"`
	MasterName            string   `json:"master-name" mapstructure:"master-name"`
	MaxIdle               int      `json:"optimisation-max-idle"    mapstructure:"optimisation-max-idle"`
	MaxActive             int      `json:"optimisation-max-active"  mapstructure:"optimisation-max-active"`
	Timeout               int      `json:"timeout"                  mapstructure:"timeout"`
	EnableCluster         bool     `json:"enable-cluster"           mapstructure:"enable-cluster"`           // 是否为集群
	UseSSL                bool     `json:"use-ssl"                  mapstructure:"use-ssl"`                  // 是否使用ssl连接
	SSLInsecureSkipVerify bool     `json:"ssl-insecure-skip-verify" mapstructure:"ssl-insecure-skip-verify"` // 是否SSL不安全跳过验证
}

func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Host:                  "127.0.0.1:6379",
		Port:                  6379, // 若 port > 0 进行替换
		Addrs:                 []string{},
		Username:              "",
		Password:              "",
		Database:              0,
		MasterName:            "",
		MaxIdle:               2000,
		MaxActive:             4000,
		Timeout:               0,
		EnableCluster:         false,
		UseSSL:                false,
		SSLInsecureSkipVerify: false,
	}
}

// AddFlags Redis相关flag
func (option *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&option.Host, "redis.host", option.Host, "flag redis host")
	fs.IntVar(&option.Port, "redis.port", option.Port, "flag redis port")
	fs.StringVar(&option.Username, "redis.username", option.Host, "flag redis username")
	fs.StringVar(&option.Password, "redis.password", option.Host, "flag redis password")
	fs.IntVar(&option.Database, "redis.database", option.Database, "flag redis database")
	fs.StringSliceVar(&option.Addrs, "redis.addrs", option.Addrs, "A set of redis address(format: 127.0.0.1:6379).")

	// todo 表示
	fs.StringVar(&option.MasterName, "redis.master-name", option.Host, "flag redis master-name")
	fs.IntVar(&option.MaxIdle, "redis.optimisation-max-idle", option.MaxIdle, "flag redis optimisation-max-idle")
	fs.IntVar(&option.MaxActive, "redis.optimisation-max-active", option.MaxIdle, "flag redis optimisation-max-active")
	fs.IntVar(&option.Timeout, "redis.timeout", option.Timeout, "flag redis timeout")

	fs.BoolVar(&option.EnableCluster, "redis.enable-cluster", option.EnableCluster, "flag redis enable-cluster") // 集群开关
	fs.BoolVar(&option.UseSSL, "redis.use-ssl", option.UseSSL, "flag redis use-ssl")
	fs.BoolVar(&option.SSLInsecureSkipVerify, "redis.ssl-insecure-skip-verify", option.SSLInsecureSkipVerify, "flag redis ssl-insecure-skip-verify")

}

// Validate 验证
func (option *RedisOptions) Validate() []error {
	return []error{}
}
