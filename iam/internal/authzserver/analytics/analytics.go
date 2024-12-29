package analytics

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"iam/pkg/cache"
	"sync"
	"sync/atomic"
	"time"
)

// Analytics将把分析数据记录到Config对象中定义的redis后端

// 分析操作 --> chan[] --> 批量写入到redis中

// 记录认证日志 -->

const analyticsKey = "iam-authorization-analytics"

var analytics *Analytics

// AnalyticsRecord 写入记录
type AnalyticsRecord struct {
	CreateTime int64     `json:"create_time"`
	Username   string    `json:"username"`
	Effect     string    `json:"effect"`
	Conclusion string    `json:"conclusion"`
	Request    string    `json:"request"`
	Policies   string    `json:"policies"`
	Deciders   string    `json:"deciders"`
	ExpireAt   time.Time `json:"expireAt"`
}

// Analytics 匹配记录认证是否通过的记录
type Analytics struct {
	store                 cache.AnalyticsHandler // 操作redis的一些内容
	poolSize              int                    // 工作线程
	recordsChan           chan *AnalyticsRecord  // 缓冲chan
	workBufferSize        int                    // buffer 个数 -> redis
	maxSyncTime           int                    // 最大多少ms进行同步
	storageExpirationTime time.Duration          // 过期时间
	stop                  uint32                 // chan关闭
	poolWg                sync.WaitGroup
}

func (a *Analytics) SetExpiration(expire int) {
	if expire == 0 {
		expire = 24
	}
	a.storageExpirationTime = time.Hour * time.Duration(expire)
}

func NewAnalytics(opt AnalyticsOptions) *Analytics {

	analytics = &Analytics{
		recordsChan:           make(chan *AnalyticsRecord, opt.RecordsBufferSize),
		poolSize:              opt.PoolSize,
		workBufferSize:        opt.RecordsBufferSize / opt.PoolSize,
		maxSyncTime:           opt.MaxSyncTime,
		storageExpirationTime: opt.StorageExpirationTime,
	}
	return analytics
}

func (a *Analytics) SetStore(store cache.AnalyticsHandler) {
	a.store = store
}

// GetAnalytics 获取全局消费日志结构
func GetAnalytics() *Analytics {
	return analytics
}

func (a *Analytics) Start() {
	// 判断redis此时是否能够正常连接

	atomic.SwapUint32(&a.stop, 0)

	for i := 0; i < a.poolSize; i++ {
		go a.workConsumption()
	}

}

func (a *Analytics) Stop() {
	atomic.SwapUint32(&a.stop, 1)
	close(a.recordsChan)
	a.poolWg.Wait()
}

func (a *Analytics) SendRecord(record *AnalyticsRecord) error {
	// 判断chan是否关闭
	if atomic.LoadUint32(&a.stop) > 0 {
		return nil
	}
	a.recordsChan <- record
	return nil
}

// 根据参数进行创建消费
func (a *Analytics) workConsumption() {

	lastSentTS := time.Now()
	buffers := make([][]byte, 0)

	a.poolWg.Add(1)
	defer a.poolWg.Done()

	for {

		var readyToSend bool
		select {

		case record, ok := <-a.recordsChan:
			// 说明改chan已经关闭，需要将该 buffers 进行存储下
			if !ok {
				a.store.AppendAnalytics(analyticsKey, buffers)
				break
			}

			bytes, err := json.Marshal(record)
			if err != nil {
				logrus.Errorf("marshal record err:%v", err)

			} else {
				buffers = append(buffers, bytes)
			}

			if len(buffers) >= a.workBufferSize {
				readyToSend = true
			}

		case <-time.After(time.Duration(a.maxSyncTime) * time.Millisecond):
			readyToSend = true
		}

		// 个数>xx个 || 时间超过多少
		timeSub := time.Now().Sub(lastSentTS).Milliseconds()

		if len(buffers) > 0 && (readyToSend || timeSub >= int64(a.maxSyncTime)) {
			a.store.AppendAnalytics(analyticsKey, buffers)
			buffers = buffers[:0]
			lastSentTS = time.Now()
		}

	}

}

/*
              创建用户
   .keep-server 创建权限  ----> 同步 redis   ----> authorization (认证该用户存在对应权限)
              创建密钥
*/
