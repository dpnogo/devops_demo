package apiserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	userv1 "iam/internal/apiserver/controller/v1/user"
	"iam/internal/apiserver/store"
	"iam/internal/pkg/middleware/auth"
	"iam/pkg/core"
	"net/http"
)

func initRouter(engine *gin.Engine) {
	installController(engine)
}

func installController(g *gin.Engine) *gin.Engine {

	strategy := auth.NewJWTStrategy(auth.NewGinGwt())

	g.POST("/login", strategy.LoginHandler)     // 登录
	g.POST("/logout", strategy.LogoutHandler)   // 登出
	g.POST("/refresh", strategy.RefreshHandler) // 刷新

	auto := newAuto()
	// 若无以下接口
	g.NoRoute(auto.Auth(), func(c *gin.Context) {
		core.WriteResponse(c, http.StatusBadRequest, fmt.Errorf("auth failed"), nil) // todo err 构建
	})

	// 获取mysql的信息
	storeIns := store.GetFactory()

	v1 := g.Group("/v1")
	{
		user := v1.Group("/user") // auto.Auth()
		userCtl := userv1.NewUserCtl(storeIns)
		user.POST("/create", userCtl.Create) // 创建用户 -->
	}

	return g
}
