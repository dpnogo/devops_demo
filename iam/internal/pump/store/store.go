package store

// AnalyticsStore 从 redis 中获取值
type AnalyticsStore interface {
	New() AnalyticsStore
	Init(cfg interface{})
	GetName() string
	GetKeysAndDel() (val []any, err error) // 从redis中获取所有值，然后进行
}

// 从 store 中获取往 pumpus 里面进行推送并写入
