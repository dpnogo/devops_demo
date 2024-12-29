package middleware

import "github.com/gin-gonic/gin"

const UsernameKey = "username"

// Context 公共使用string值  --> 暂时
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(UsernameKey, c.GetString(UsernameKey))
		c.Next()
	}
}
