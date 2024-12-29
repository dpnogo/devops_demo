package auth

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"strings"
)

type BasicStrategy struct {
	verify func(username, password string) bool // 验证账号密码
}

func NewBasic(verify func(username, password string) bool) *BasicStrategy {
	return &BasicStrategy{
		verify,
	}
}

// Auth Authorization: Basic ${basic}
func (basic *BasicStrategy) Auth() gin.HandlerFunc {

	return func(c *gin.Context) {

		basicStr := c.Request.Header.Get("Authorization")

		if basicStr == "" {
			// code basic is empty

			c.Abort()
			return
		}

		// Basic username:password

		auth := strings.SplitN(basicStr, " ", 2)
		if len(auth) != 2 {
			// basic len not inconformity

			c.Abort()
			return
		}

		if auth[0] != "Basic" {
			// auth type is not basic

			c.Abort()
			return
		}

		bytes, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			// base auth decode failed

			c.Abort()
			return
		}

		infoStr := string(bytes)
		info := strings.Split(infoStr, ":")

		if len(info) != 2 || basic.verify(info[0], info[1]) {
			// code 构建 info len or username/password verification failure

			c.Abort()
			return
		}

		// 返回对应的信息
		c.Next()

	}

}
