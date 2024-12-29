package authzserver

import (
	"github.com/gin-gonic/gin"
	"iam/internal/authzserver/controller/v1/authorize"
	"iam/internal/authzserver/load/cache"
	"log"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
	return
}

func installController(g *gin.Engine) {

	// 认证身份
	//auth := newCacheAuth()
	//g.NoRoute(auth.AuthFunc(), func(c *gin.Context) {
	//	core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "page not found."), nil)
	//})

	cacheIns, _ := cache.GetCacheInsOr(nil)
	if cacheIns == nil {
		log.Panicf("get nil cache instance")
	}

	apiv1 := g.Group("/v1") //  auth.AuthFunc()
	{
		authzController := authorize.NewAuthorizeCtl(cacheIns)

		// Router for authorization
		apiv1.POST("/authorization", authzController.Authorize)
	}

}
