package v1

import (
	"context"
	"iam/internal/apiserver/store"
	"iam/pkg/api/user"
)

type UserSvc interface {
	CreateUser(ctx context.Context, user *user.User) error
	DeleteUser(ctx context.Context, userId uint64) error
	// 批量删除
	UpdateUser(ctx context.Context, user *user.User) error
	GetUser(ctx context.Context, userId uint64) (*user.User, error)
}

type userSvc struct {
	factory store.Factory
}

func newUserSvc(f store.Factory) *userSvc {
	return &userSvc{f}
}

func (svc *userSvc) CreateUser(ctx context.Context, user *user.User) error {
	return svc.factory.User().CreateUser(ctx, user)
}

func (svc *userSvc) DeleteUser(ctx context.Context, userId uint64) error {
	return svc.factory.User().DeleteUser(ctx, userId)
}

func (svc *userSvc) UpdateUser(ctx context.Context, user *user.User) error {
	return svc.factory.User().UpdateUser(ctx, user)
}

func (svc *userSvc) GetUser(ctx context.Context, userId uint64) (*user.User, error) {
	return svc.factory.User().GetUser(ctx, userId)
}
