package cache

type AnalyticsHandler interface {
	Connect() bool                              // 判断是否能够连接
	AppendAnalytics(key string, bytes [][]byte) // 存储写入到redis的值
}
