package options

import (
	"github.com/spf13/pflag"
	"time"
)

type MysqlOptions struct {
	Host     string `json:"host" mapstructure:"host"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Database string `json:"database"  mapstructure:"database"`

	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"` // mysql 最大空闲连接数
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"` // mysql 最大打开连接数
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int           `json:"log-level"                          mapstructure:"log-level"`
}

func NewMysqlOptions() *MysqlOptions {
	return &MysqlOptions{
		Host:                  "127.0.0.1:3306",
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              1, // Silent
	}
}

func (option *MysqlOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&option.Host, "mysql.host", option.Host, "mysql host")
	fs.StringVar(&option.Username, "mysql.username", option.Username, "mysql username")
	fs.StringVar(&option.Password, "mysql.password", option.Password, "mysql password")
	fs.StringVar(&option.Database, "mysql.database", option.Database, "mysql database")

	// todo 描述
	fs.IntVar(&option.MaxIdleConnections, "mysql.max-idle-connections", option.MaxIdleConnections, "mysql max-idle-connections")
	fs.IntVar(&option.MaxOpenConnections, "mysql.max-open-connections", option.MaxOpenConnections, "mysql max-open-connections")
	fs.DurationVar(&option.MaxConnectionLifeTime, "mysql.max-connection-life-time", option.MaxConnectionLifeTime, "mysql max-connection-life-time")
	fs.IntVar(&option.LogLevel, "mysql.log-level", option.LogLevel, "mysql log-level")
}

func (option *MysqlOptions) Validate() []error {
	return []error{}
}
