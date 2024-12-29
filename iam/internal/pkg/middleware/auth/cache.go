package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"iam/internal/pkg/middleware"
	"time"
)

var (
	ErrMissingKID    = errors.New("invalid token format: missing kid field in claims")
	ErrMissingSecret = errors.New("can not obtain secret information from cache")
)

// Secret 包含密钥的基本信息
type Secret struct {
	Username string
	ID       string
	Key      string
	Expires  int64
}

type CacheStrategy struct {
	get func(kid string) (Secret, error) // 根据kid的得到对应的密钥信息
}

func NewCacheStrategy(getSecret func(kid string) (Secret, error)) CacheStrategy {
	return CacheStrategy{getSecret}
}

func (cs *CacheStrategy) Auth() gin.HandlerFunc {

	// 从head中获取 Authorization
	return func(c *gin.Context) {

		key := c.Request.Header.Get("Authorization")
		if key == "" {
			// 返回错误信息 Authorization 为空
			c.Abort()
			return
		}

		var (
			tokenStr string
			secret   Secret
		)

		// 验证 Authorization 的 Bearer xxx 字段
		fmt.Sscanf(key, "Bearer %s", &tokenStr)

		if tokenStr == "" {
			// 返回错误信息 token 为空
			c.Abort()
			return
		}

		// 进行解析token
		claims := &jwt.MapClaims{}

		parsedT, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {

			// 验证 tokenStr 值的算法的HMAC签名
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, ErrMissingKID
			}
			var err error
			secret, err = cs.get(kid)
			if err != nil {
				return "", err
			}

			return []byte(secret.Key), nil
		})

		// 存在错误 || token无效
		if err != nil || !parsedT.Valid {
			c.Abort()
			// 构建 token 无效的msg
			return
		}

		// 验证是否过期
		if secret.Expires >= time.Now().Unix() {
			// 过期 msg
			c.Abort()
			return
		}

		c.Set(middleware.UsernameKey, secret.Username) // 设置人名

		c.Next()
	}

}
