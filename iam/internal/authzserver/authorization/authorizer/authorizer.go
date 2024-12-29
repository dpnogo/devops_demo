package authorizer

import (
	"github.com/ory/ladon"
	"iam/internal/authzserver/authorization"
)

type PolicyGetter interface {
	GetPolicy(key string) ([]*ladon.DefaultPolicy, error) // 通过 key 得到所有的策略 -->
}

// Authorization 实现 Authorization.interface
type Authorization struct {
	getter PolicyGetter
}

// NewAuthorization 初始化
func NewAuthorization(getter PolicyGetter) authorization.AuthorizationInterface {
	return &Authorization{getter}
}

// Create 创建持久化策略 --> 返回nil，因为我们使用mysql存储来存储策略
func (auth *Authorization) Create(policy *ladon.DefaultPolicy) error {
	return nil
}

// Update 更新现有策略
func (auth *Authorization) Update(policy *ladon.DefaultPolicy) error {
	return nil
}

// Get 获取当前id的策略
func (auth *Authorization) Get(id string) (*ladon.DefaultPolicy, error) {
	return nil, nil
}

// Delete 删除当前id的策略
func (auth *Authorization) Delete(id string) error {
	return nil
}

func (auth *Authorization) BatchDelete(ids []string) error {
	return nil
}

// GetAll 获取某一页的策略
func (auth *Authorization) GetAll(limit int64, offset int64) (ladon.Policies, error) {
	return nil, nil
}

// List 获取某个name的所有的,从 db(cache)中获取
func (auth *Authorization) List(username string) ([]*ladon.DefaultPolicy, error) {
	return auth.getter.GetPolicy(username)
}

// LogRejectedAccessRequest 将认证失败请求日志写到一个统一的chan中，进行后台消费
func (auth *Authorization) LogRejectedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {

}

// LogGrantedAccessRequest 将认证成功日志写到一个统一的chan中，进行后台消费
func (auth *Authorization) LogGrantedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {

}
