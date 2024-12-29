package store

// 定义上方使用的接口

// 默认 client
var client Factory

type Factory interface {
	Policies() PolicyStore
	Secrets() SecretStore
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
