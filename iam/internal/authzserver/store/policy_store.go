package store

import "github.com/ory/ladon"

// 定义策略相关方法

type PolicyStore interface {
	List() (map[string][]*ladon.DefaultPolicy, error) // 获取所有用户的所有策略
	Get(key string) ([]*ladon.DefaultPolicy, error)   // 获取某一个人/用户的策略
}
