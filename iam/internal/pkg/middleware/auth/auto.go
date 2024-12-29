package auth

import (
	"github.com/gin-gonic/gin"
	"iam/internal/pkg/middleware"
	"strings"
)

const (
	AuthorizationBasic  = "Basic"
	AuthorizationBearer = "Bearer"
)

// 根据 Authorization 头部为 Basic和Bearer 来进行分别使用对应的验证方式

// AutoStrategy 接口
type AutoStrategy struct {
	basic middleware.AuthStrategy
	jwt   middleware.AuthStrategy
}

func NewAutoStrategy(basic, jwt middleware.AuthStrategy) *AutoStrategy {
	return &AutoStrategy{
		basic,
		jwt,
	}
}

// Auth Authorization xxx
func (auto *AutoStrategy) Auth() gin.HandlerFunc {

	return func(c *gin.Context) {

		auth := c.Request.Header.Get("Authorization")
		auths := strings.Split(auth, " ")
		if len(auths) == 0 {
			// Authorization key 不符合条件

			c.Abort()
			return
		}

		authPolicy := &middleware.AuthPolicy{}
		switch auths[0] {

		case AuthorizationBasic:
			authPolicy.SetPolicy(auto.basic)

		case AuthorizationBearer:
			authPolicy.SetPolicy(auto.jwt)

		default:
			// 未支持的 Authorization 类型

			c.Abort()
			return
		}

		// 执行
		authPolicy.AuthFunc()
		c.Next()

	}

}
