package config

import "iam/internal/apiserver/options"

type Config struct {
	*options.Options
}

func NewConfig(opt *options.Options) *Config {
	return &Config{opt}
}
