package mysql

import (
	"context"
	"gorm.io/gorm"
	"iam/pkg/api/user"
)

type userStore struct {
	db *gorm.DB
}

func newUsers(ds *gorm.DB) *userStore {
	return &userStore{ds}
}

// CreateUser 添加 user
func (store *userStore) CreateUser(ctx context.Context, user *user.User) error {

	return store.db.Create(&user).Error
}

// DeleteUser 删除 user
func (store *userStore) DeleteUser(ctx context.Context, userId uint64) error {
	// 删除该用户相关的策略,
	return nil
}

// UpdateUser 更新用户信息
func (store *userStore) UpdateUser(ctx context.Context, user *user.User) error {
	return store.db.Save(user).Error
}

func (store *userStore) GetUser(ctx context.Context, userId uint64) (user *user.User, err error) {
	err = store.db.Where("id = ? and status = 1", userId).Take(user).Error
	return
}

func (store *userStore) GetUserByName(ctx context.Context, username string) (user *user.User, err error) {
	err = store.db.Where("username = ? and status = 1", username).Take(user).Error
	return
}
