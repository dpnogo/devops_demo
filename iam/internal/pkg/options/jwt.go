package options

import (
	"github.com/spf13/pflag"
	"iam/internal/pkg/server"
	"time"
)

type JwtOptions struct {
	Realm      string        `json:"realm"       mapstructure:"realm"`
	Key        string        `json:"key"         mapstructure:"key"` // 私钥
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
}

func NewJwtOptions() *JwtOptions {
	defaults := server.NewConfig()

	return &JwtOptions{
		Realm:      defaults.JwtInfo.Realm,
		Key:        defaults.JwtInfo.Key,
		Timeout:    defaults.JwtInfo.Timeout,
		MaxRefresh: defaults.JwtInfo.MaxRefresh,
	}
}

func (s *JwtOptions) ApplyTo(config *server.Config) error {
	config.JwtInfo = &server.JwtInfo{
		Realm:      s.Realm,
		Key:        s.Key,
		Timeout:    s.Timeout,
		MaxRefresh: s.MaxRefresh,
	}
	return nil
}

func (s *JwtOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Realm, "jwt.realm", s.Realm, "Realm name to display to the user.")
	fs.StringVar(&s.Key, "jwt.key", s.Key, "用于签署jwt令牌的私钥.")
	fs.DurationVar(&s.Timeout, "jwt.timeout", s.Timeout, "JWT token timeout.")
	fs.DurationVar(&s.MaxRefresh, "jwt.max-refresh", s.MaxRefresh,
		"This field allows clients to refresh their token until MaxRefresh has passed.") // 这个字段允许客户端刷新他们的令牌，直到MaxRefresh通过。

}

func (s *JwtOptions) Validate() []error {
	return []error{}
}
