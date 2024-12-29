package shutdown

import "sync"

type ShutdownCallback interface {
	OnShutdown(string) error
}
type ShutdownFunc func(string) error

// OnShutdown defines the action needed to run when shutdown triggered.
func (f ShutdownFunc) OnShutdown(shutdownManager string) error {
	return f(shutdownManager)
}

type ShutdownManager interface {
	GetName() string
	Start(gs GSInterface) error
	ShutdownStart() error
	ShutdownFinish() error
}

// ErrorHandler 是一个接口，你可以传递给SetErrorHandler来处理异步错误。
type ErrorHandler interface {
	OnError(err error)
}

// GSInterface 是一个由GracefulShutdown实现的接口
// 它被传递给ShutdownManager在关机被请求时调用StartShutdown。
type GSInterface interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(shutdownCallback ShutdownCallback)
}

// GracefulShutdown 是处理 ShutdownCallbacks 和 ShutdownManagers 的主要结构体。用New初始化它。
type GracefulShutdown struct {
	callbacks    []ShutdownCallback // 回调函数
	managers     []ShutdownManager  // 管理Shutdown
	errorHandler ErrorHandler       // 回调函数时候，出现err后执行函数
}

// New initializes GracefulShutdown.
func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]ShutdownCallback, 0, 10),
		managers:  make([]ShutdownManager, 0, 3),
	}
}

// AddShutdownManager 添加了一个ShutdownManager，它将监听关机请求。
func (gs *GracefulShutdown) AddShutdownManager(manager ShutdownManager) {
	gs.managers = append(gs.managers, manager)
}

//type GSInterface interface {
//	StartShutdown(sm ShutdownManager)
//	ReportError(err error)
//	AddShutdownCallback(shutdownCallback ShutdownCallback)
//}

func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.managers {
		// gs 要实现 GSInterface 接口
		if err := manager.Start(gs); err != nil {
			return err
		}
	}
	return nil
}

// GracefulShutdown 实现 GSInterface 接口

func (gs *GracefulShutdown) AddShutdownCallback(shutdownCallback ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, shutdownCallback)
}

// StartShutdown 信号发生后，进行执行回调函数
func (gs *GracefulShutdown) StartShutdown(sm ShutdownManager) {
	gs.ReportError(sm.ShutdownStart())

	var wg sync.WaitGroup
	for _, shutdownCallback := range gs.callbacks {
		wg.Add(1)
		go func(shutdownCallback ShutdownCallback) {
			defer wg.Done()

			gs.ReportError(shutdownCallback.OnShutdown(sm.GetName()))
		}(shutdownCallback)
	}

	wg.Wait()

	gs.ReportError(sm.ShutdownFinish())
}

func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}
