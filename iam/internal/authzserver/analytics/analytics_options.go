package analytics

import (
	"fmt"
	"github.com/spf13/pflag"
	"time"
)

// AnalyticsOptions 配置
type AnalyticsOptions struct {
	PoolSize                int           `json:"pool-size" mapstructure:"pool-size"`                                 // 工作线程个数
	MaxSyncTime             int           `json:"max-sync-time" mapstructure:"max-sync-time"`                         // 最大多少ms进行同步,时间间隔
	RecordsBufferSize       int           `json:"records-buffer-size" mapstructure:"records-buffer-size"`             // 最大缓冲区大小 / 线程 = 每个线程1执行最大大小
	StorageExpirationTime   time.Duration `json:"storage-expiration-time"   mapstructure:"storage-expiration-time"`   // 过期时间
	Enable                  bool          `json:"enable"                    mapstructure:"enable"`                    // 是否启用缓冲日志
	EnableDetailedRecording bool          `json:"enable-detailed-recording" mapstructure:"enable-detailed-recording"` // 启用详细记录(暂时)
}

func NewAnalyticsOptions() *AnalyticsOptions {
	return &AnalyticsOptions{
		PoolSize:              3,
		MaxSyncTime:           500,
		RecordsBufferSize:     300,
		StorageExpirationTime: time.Hour * 24,
	}
}

func (option *AnalyticsOptions) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&option.PoolSize, "analytics.pool-size", option.PoolSize, "analytics pool-size")
	fs.IntVar(&option.MaxSyncTime, "analytics.max-sync-time", option.MaxSyncTime, "analytics max-sync-time")
	fs.IntVar(&option.RecordsBufferSize, "analytics.records-buffer-size", option.RecordsBufferSize, "analytics records-buffer-size")

	fs.DurationVar(&option.StorageExpirationTime, "analytics.storage-expiration-time", option.StorageExpirationTime,
		"analytics storage-expiration-time")

	fs.BoolVar(&option.Enable, "analytics.enable", option.Enable, "analytics enable")
	fs.BoolVar(&option.EnableDetailedRecording, "analytics.enable-detailed-recording", option.EnableDetailedRecording,
		"analytics enable-detailed-recording")
}

func (option *AnalyticsOptions) Validate() []error {
	if option == nil {
		return nil
	}
	errors := []error{}

	if option.Enable && (option.MaxSyncTime < 1 || option.MaxSyncTime > 1000) {
		errors = append(errors, fmt.Errorf("--analytics.flush-interval %v must be between 1 and 1000", option.MaxSyncTime))
	}

	return errors
}
