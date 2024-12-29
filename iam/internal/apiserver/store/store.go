package store

// 定义一个全局使用的client

var client Factory

type Factory interface {
	User() UserStore
	Close() error
}

func SetFactory(factory Factory) {
	client = factory
}

func GetFactory() Factory {
	return client
}
