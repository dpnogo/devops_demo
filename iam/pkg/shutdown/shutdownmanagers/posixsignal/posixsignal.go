package posixsignal

import (
	"iam/pkg/shutdown"
	"os"
	"os/signal"
	"syscall"
)

// Name 定义 shutdown manager 名称
const Name = "PosixSignalManager"

// PosixSignalManager 实现了添加到 GracefulShutdown 的 ShutdownManager 接口。
// 用NewPosixSignalManager初始化。
type PosixSignalManager struct {
	signals []os.Signal // 信号
}

/*
type ShutdownManager interface {
    GetName() string
    Start(gs GSInterface) error
    ShutdownStart() error
    ShutdownFinish() error
}
*/

func (p *PosixSignalManager) GetName() string {
	return Name
}

// Start 等待信号输入
func (p *PosixSignalManager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, p.signals...)
		// Block until a signal is received.
		<-c
		// 若存在对应信号，进行开始shutdown，即执行回调函数
		gs.StartShutdown(p)
	}()

	return nil
}

func (p *PosixSignalManager) ShutdownStart() error {
	return nil
}

func (p *PosixSignalManager) ShutdownFinish() error {
	os.Exit(0)

	return nil
}

// NewPosixSignalManager 初始化PosixSignalManager。
// 作为参数，你可以提供os。要侦听的信号s，如果没有给出，则默认为SIGINT和SIGTERM。
func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}

	return &PosixSignalManager{
		signals: sig,
	}
}
