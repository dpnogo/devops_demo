package config

import "iam/internal/authzserver/options"

// 供外部使用的配置

type Config struct {
	*options.Options
}

// NewConfigFromOption 根据构建的option创建config
func NewConfigFromOption(o *options.Options) *Config {
	return &Config{o}
}
