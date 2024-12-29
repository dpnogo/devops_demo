package v1

import "iam/internal/apiserver/store"

// store 仅仅做和db操作内容，这里进行相关的逻辑操作

type Service interface {
	User() UserSvc
}

type service struct {
	factory store.Factory
}

func (svc *service) User() UserSvc {
	return newUserSvc(svc.factory)
}

// NewSvc 外部使用服务，返回对应操作的接口
func NewSvc(factory store.Factory) Service {
	return &service{factory}
}
