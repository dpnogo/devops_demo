package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
)

const (
	// XRequestIDKey 定义X-Request-ID值
	XRequestIDKey = "X-Request-ID"
)

// RequestId 为每一系列请求创建一个id，并set到head中
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {

		rid := c.GetHeader(XRequestIDKey)

		if rid == "" {
			var err error
			rid = uuid.Must(uuid.NewV4(), err).String()
			if err != nil {
				log.Printf("create request init uuid err:%s\n", err)
			}
			c.Request.Header.Set(XRequestIDKey, rid)
			c.Set(XRequestIDKey, rid)
		}

	}
}
