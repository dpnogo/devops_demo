package middleware

import (
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
	"net/http"
	"time"
)

var Middlewares = defaultMiddlewares()

func Secure() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 访问控制允许来源 --> * 浏览器允许来自任何来源的代码访问资源
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")

		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000") // 设置 HSTS 头部，所有现有和未来的子域名将在 1 年内使用 HTTPS
		}
	}
}

/*
 http header:
      Access-Control-Allow-Origin: https://developer.mozilla.org  允许 https://developer.mozilla.org 的请求访问资源
      X-Frame-Options: "DENY","SAMEORIGIN","ALLOW-FROM origin"   浏览器是否允许在<frame>、<iframe>或<object>中呈现页面
           DENY: 无论从哪个站点加载页面，浏览器都会禁止在框架中加载页面
           SAMEORIGIN: 只有与页面本身相同来源的站点才能加载页面
           ALLOW-FROM origin: 允许页面仅在指定的来源和/或域中加载
      X-Content-Type-Options :  用于指示服务器广告中的Content-Type头应该被遵循而不被更改,这个头允许您避免MIME类型嗅探，即告诉浏览器MIME类型是经过故意配置的。
      X-XSS-Protection: 启用现代Web浏览器中内置的跨站脚本（XSS）过滤器。通常情况下，这个过滤器是默认启用的
      Strict-Transport-Security: 告知浏览器只能通过 HTTPS 访问网站，并且任何尝试使用 HTTP 访问的请求都应自动转换为 HTTPS 请求。
*/

// Options 检测服务器所支持的请求方法 + CORS 中的预检请求(preflight request)
func Options() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Content-Type", "application/json")
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.

func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		c.Next()
	}
}

func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery":  gin.Recovery(),
		"secure":    Secure(), // https
		"options":   Options(),
		"nocache":   NoCache(),
		"cors":      Cors(),
		"requestid": RequestId(),
		//"logger":    Logger(),
		"dump": gindump.Dump(), // 中间件来打印请求和响应的头部和主体
	}
}
