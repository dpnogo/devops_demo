package load

import (
	"context"
	"log"
	"sync"
)

/*

 auth (更新等操作) 新增/删除/修改 --> 中间件 ()    --> todo 更新缓存/更新db
   \
    V 是否有修改 --> 批量操作(例如), 去掉某个用户的某些策略-->
    redis   ----->  authorization  --> 更新到内存库中

   // 首次加载 xx 个， 后续新增/删除/更新






   authorization --> 操作到内存中()

*/

// Loader 将密钥和策略添加到被内存中
type Loader interface {
	Reload() error // 全量加载 (或者这里进行限制大小)
	// 其他单值操作 -->
}

type Load struct {
	ctx    context.Context
	lock   *sync.RWMutex
	loader Loader
}

// NewLoader 初始化
func NewLoader(ctx context.Context, loader Loader) *Load {
	return &Load{
		ctx:    ctx,
		lock:   new(sync.RWMutex),
		loader: loader,
	}
}

func (l *Load) Start() {

	// 订阅redis key进行后台加载到缓存中
	go startPubSubLoop()

	// todo 全量加载，进行优化

	// 刚开始时候全量加载
	l.DoReload()
}

func startPubSubLoop() {

}

func (l *Load) DoReload() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if err := l.loader.Reload(); err != nil {
		log.Printf("faild to refresh target storage: %s\n", err.Error())
	}

	log.Println("refresh target storage succ")
}
