package pumps

import "iam/internal/pump/analytics"

// Pump 分析接口，抽象成接口 , 支持不同上报服务，可提供以插件的方式
type Pump interface {
	GetName() string
	New() Pump
	Init(interface{}) error      //
	WriteData(datas []any) error // 往下游系统写入数据

	SetFilters(analytics.AnalyticsFilters) // 设置是否过滤某条数据
	GetFilters() analytics.AnalyticsFilters
	SetTimeout(timeout int) // 设置超时时间
	GetTimeout() int
	SetOmitDetailedRecording(bool)  // 过滤掉详细的数据
	GetOmitDetailedRecording() bool // 是否过滤详情数据
}
