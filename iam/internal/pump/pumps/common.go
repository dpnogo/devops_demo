package pumps

import "iam/internal/pump/analytics"

// CommonPumpConfig 通用选项配置
type CommonPumpConfig struct {
	filters               analytics.AnalyticsFilters
	timeout               int
	omitDetailedRecording bool
}

func (comConfig *CommonPumpConfig) SetFilters(filters analytics.AnalyticsFilters) {
	comConfig.filters = filters
}

func (comConfig *CommonPumpConfig) GetFilters() analytics.AnalyticsFilters {
	return comConfig.filters
}

func (comConfig *CommonPumpConfig) SetTimeout(timeout int) {
	comConfig.timeout = timeout
}

func (comConfig *CommonPumpConfig) GetTimeout() int {
	return comConfig.timeout
}

func (comConfig *CommonPumpConfig) SetOmitDetailedRecording(odr bool) {
	comConfig.omitDetailedRecording = odr
}

func (comConfig *CommonPumpConfig) GetOmitDetailedRecording() bool {
	return comConfig.omitDetailedRecording
}
