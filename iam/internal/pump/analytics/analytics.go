package analytics

// AnalyticsFilters 定义分析选项
type AnalyticsFilters struct {
	Usernames        []string `json:"usernames"`
	SkippedUsernames []string `json:"skip_usernames"`
}
