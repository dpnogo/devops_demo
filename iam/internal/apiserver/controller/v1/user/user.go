package user

import (
	svcv1 "iam/internal/apiserver/service/v1"
	"iam/internal/apiserver/store"
)

type UserController struct {
	svc svcv1.Service
}

func NewUserCtl(factory store.Factory) *UserController {
	return &UserController{svc: svcv1.NewSvc(factory)}
}
