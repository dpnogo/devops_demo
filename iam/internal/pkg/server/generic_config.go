package server

import (
	"github.com/gin-gonic/gin"
	"time"
)

/*
  pkg.Options --> 选择后组成 应用配置 --> 选择服务相关配置(构建服务器)
*/

// 应用配置 ApplyTo  --> 服务器配置

// Config HTTP 服务器的配置,启动服务器所需的配置
type Config struct {
	InsecureServing *InsecureServingInfo // http
	SecureServing   *SecureServing       // https
	JwtInfo         *JwtInfo             // jwt info

	Mode        string
	Middlewares []string
	Healthz     bool // 是否添加检查

	EnableProfiling bool // 开启分析 即 pprof
	EnableMetrics   bool // 开启 metrics
}

// InsecureServingInfo 不安全服务选项
type InsecureServingInfo struct {
	Address string
}

// SecureServing 安全服务选项
type SecureServing struct {
	Address string  `json:"address"` // bind_address:bind_port
	CertKey CertKey `json:"tls"`
}

type CertKey struct {
	// CertFile 文件是否包含pem编码的证书，可能还包含完整的证书链
	CertFile string
	// KeyFile 文件是否包含pem编码的私钥，用于CertFile指定的证书
	KeyFile string
}

type JwtInfo struct {
	// defaults to "iam jwt"
	Realm string
	// defaults to empty
	Key string //
	// defaults to one hour
	Timeout time.Duration // token 超时时间
	// defaults to zero
	MaxRefresh time.Duration // 最大刷新
}

// NewConfig 默认config配置
func NewConfig() *Config {
	return &Config{
		Healthz:         false,
		Mode:            gin.ReleaseMode, // 运行模式: "debug","test","release"
		Middlewares:     []string{},
		EnableProfiling: true,
		EnableMetrics:   true,
		JwtInfo: &JwtInfo{
			Realm:      "iam jwt",
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		},
	}
}

// CompletedConfig 是GenericAPIServer的完整配置，提供给外部使用
type CompletedConfig struct {
	*Config
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

// New 通过 config 初始化 GenericAPIServer ,  todo 不直接使用config构建? 而是在包一层? 说明已经应用过配置
func (c CompletedConfig) New() (*GenericAPIServer, error) {

	gin.SetMode(c.Mode)

	s := &GenericAPIServer{
		SecureServing:   c.SecureServing,
		InsecureServing: c.InsecureServing,
		healthz:         c.Healthz,
		enableMetrics:   c.EnableMetrics,
		enableProfiling: c.EnableProfiling,
		middlewares:     c.Middlewares,
		Engine:          gin.New(),
	}
	// 根据选项加载配置所需功能
	initGenericAPIServer(s)

	return s, nil
}
