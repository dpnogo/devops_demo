package options

import (
	"github.com/spf13/pflag"
	"iam/internal/pkg/server"
)

// FeatureOptions 功能 是否开启相关分析功能
type FeatureOptions struct {
	FeatureProfiling bool `json:"profiling" mapstructure:"profiling"`           // 开启分析 即 pprof
	FeatureMetrics   bool `json:"enable-metrics" mapstructure:"enable-metrics"` // 开启 metrics
}

func NewFeatureOptions() *FeatureOptions {
	defaultCfg := server.NewConfig()
	return &FeatureOptions{
		FeatureProfiling: defaultCfg.EnableProfiling,
		FeatureMetrics:   defaultCfg.EnableMetrics,
	}
}

func (feature *FeatureOptions) AddFlags(fs *pflag.FlagSet) {

	fs.BoolVar(&feature.FeatureMetrics, "metrics", feature.FeatureMetrics, "Enables metrics on the apiserver at /metrics")
	fs.BoolVar(&feature.FeatureProfiling, "profiling", feature.FeatureProfiling, "Enable profiling via web interface host:port/debug/pprof/")

}

func (feature *FeatureOptions) ApplyTo(config *server.Config) error {
	config.EnableMetrics = feature.FeatureMetrics
	config.EnableProfiling = feature.FeatureProfiling
	return nil
}

// Validate 验证各个参数
func (feature *FeatureOptions) Validate() []error {
	return []error{}
}
