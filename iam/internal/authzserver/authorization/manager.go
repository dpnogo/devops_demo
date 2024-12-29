package authorization

import "github.com/ory/ladon"

// PolicyManager 策略管理器
type PolicyManager struct {
	client AuthorizationInterface // 使用授权接口，进行对manager进行授权/去除等操作
}

func NewManager(client AuthorizationInterface) ladon.Manager {
	return &PolicyManager{
		client: client,
	}
}

// Create 创建持久化策略
func (m *PolicyManager) Create(policy ladon.Policy) error {
	return nil
}

// Update 更新现有策略
func (m *PolicyManager) Update(policy ladon.Policy) error {
	return nil
}

// Get 获取当前id的策略
func (m *PolicyManager) Get(id string) (ladon.Policy, error) {
	return nil, nil
}

// Delete 删除当前id的策略
func (m *PolicyManager) Delete(id string) error {
	return nil
}

// GetAll 获取某一页的策略
func (m *PolicyManager) GetAll(limit int64, offset int64) (ladon.Policies, error) {
	return nil, nil
}

// FindRequestCandidates 返回与请求对象匹配的候选对象。它要么返回与请求完全匹配的集合，要么返回它的超集。如果发生错误，它将返回nil和错误。 (自定义匹配)
func (m *PolicyManager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {

	username := ""

	if user, ok := r.Context["username"].(string); ok {
		username = user
	}

	// 查询某个用户的全部权限
	policies, err := m.client.List(username)
	if err != nil {
		return nil, err
	}

	pls := make([]ladon.Policy, 0)
	for _, policy := range policies {
		pls = append(pls, policy)
	}

	return pls, nil
}

// FindPoliciesForSubject 根据主题获取所有的权限
func (m *PolicyManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	return nil, nil
}

// FindPoliciesForResource 根据访问资源对应获取对应的权限
func (m PolicyManager) FindPoliciesForResource(resource string) (ladon.Policies, error) {
	return nil, nil
}
