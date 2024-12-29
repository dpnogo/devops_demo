package authorization

import (
	"github.com/ory/ladon"
	authzv1 "iam/pkg/api/authz/v1"
)

// Authorizer 授权人，执行 IsAllowed(r *Request) error，
type Authorizer struct {
	warden ladon.Warden // 用于执行访问控制策略 校验是否通过  --> 用于根据访问请求和策略进行访问控制决策，并返回相应的访问控制结果。
}

func NewAuthorizer(authorizationClient AuthorizationInterface) *Authorizer {
	return &Authorizer{
		warden: &ladon.Ladon{
			Manager:     NewManager(authorizationClient),     // 管理器
			AuditLogger: NewAuditLogger(authorizationClient), // 策略日志
		},
	}
}

// 创建/删除/
// 查找:list

func (at *Authorizer) Authorize(request *ladon.Request) *authzv1.Response {

	// todo  IsAllowed 验证的关键
	err := at.warden.IsAllowed(request)
	if err != nil {
		return &authzv1.Response{
			Denied: true,
			Reason: err.Error(),
		}
	}

	// IsAllowed 参数调用
	// l.Manager.FindRequestCandidates(r) 返回 policies ,即ladon对应的Manager , 存在错误则请求--> 进行记录 l.metric().RequestProcessingError(*r, nil, err)
	// 执行 DoPoliciesAllow --> Request和[]Policy 对于给定的策略列表，如果主题s对资源r具有上下文c的权限p，则DoPoliciesAllow返回nil，否则返回错误。

	return &authzv1.Response{
		Allowed: true,
	}
}
