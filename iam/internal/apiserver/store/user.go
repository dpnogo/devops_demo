package store

import (
	"context"
	"iam/pkg/api/user"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *user.User) error
	DeleteUser(ctx context.Context, userId uint64) error
	// 批量删除
	UpdateUser(ctx context.Context, user *user.User) error
	GetUser(ctx context.Context, userId uint64) (*user.User, error)
	GetUserByName(ctx context.Context, username string) (*user.User, error)
}
