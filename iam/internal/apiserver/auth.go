package apiserver

import (
	"context"
	"github.com/sirupsen/logrus"
	"iam/internal/apiserver/store"
	"iam/internal/pkg/middleware"
	"iam/internal/pkg/middleware/auth"
)

func newAuto() middleware.AuthStrategy {
	return auth.NewAutoStrategy(newBasic(), newJwt())
}

func newBasic() middleware.AuthStrategy {

	return auth.NewBasic(func(username, password string) bool {
		// 验证用户密码是否正确

		userInfo, err := store.GetFactory().User().GetUserByName(context.TODO(), username)
		if err != nil {
			logrus.Errorf("search username:%s, err:%v", username, err)
			return false
		}

		if err = userInfo.Compare(password); err != nil {
			logrus.Errorf("compare username:%s pwd err:%v", username, err)
			return false
		}

		return true
	})
}

func newJwt() middleware.AuthStrategy {
	return auth.NewJWTStrategy(auth.NewGinGwt())
}
