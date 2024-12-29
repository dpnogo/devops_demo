package auth

import (
	"encoding/base64"
	"fmt"
	ginJwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"iam/internal/apiserver/store"
	"iam/internal/pkg/middleware"
	"iam/pkg/api/user"
	"net/http"
	"strings"
	"time"
)

const (
	Realm           = ""
	APIServerIssuer = ".keep-server"
)

// 登录所需
type loginInfo struct {
	Username string `form:"username" json:"username" binding:"required,username"`
	Password string `form:"password" json:"password" binding:"required,password"`
}

type JWTStrategy struct {
	ginJwt.GinJWTMiddleware
}

func (jwt *JWTStrategy) Auth() gin.HandlerFunc {
	return jwt.MiddlewareFunc()
}

// NewJWTStrategy 初始化 gin jwt 验证
func NewJWTStrategy(gjwt *ginJwt.GinJWTMiddleware) *JWTStrategy {
	return &JWTStrategy{*gjwt}
}

// 登录，登出，重新刷新token

func NewGinGwt() *ginJwt.GinJWTMiddleware {

	// 实现认证函数
	gJwt, _ := ginJwt.New(&ginJwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.realm"),         // JWT # jwt 标识
		SigningAlgorithm: "HS256",                              // 签名算法
		Key:              []byte(viper.GetString("jwt.key")),   // 服务端密钥
		Timeout:          viper.GetDuration("jwt.timeout"),     // 过期
		MaxRefresh:       viper.GetDuration("jwt.max_refresh"), // 最大重试
		Authenticator:    authenticator(),                      // 登录
		PayloadFunc:      payloadFunc(),                        // 负载
		LoginResponse:    loginResponse(),                      // 登录返回
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, "ok")
		},
		RefreshResponse: refreshResponse(), // 重新登录返回值
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := ginJwt.ExtractClaims(c)
			return claims[ginJwt.IdentityKey]
		},
		IdentityKey:  middleware.UsernameKey,
		Authorizator: authorizator(), // ?
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	return gJwt
}

// data 是用户身份验证时使用 PayloadFunc 函数返回的数据
func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		//if _, ok := data.(string); ok {
		//
		//	return true
		//}

		fmt.Println("data", data)
		return true
	}
}

// 重新登录，返回值
func refreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

// 登录成功后返回
func loginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func payloadFunc() func(data interface{}) ginJwt.MapClaims {

	return func(data interface{}) ginJwt.MapClaims {
		claims := ginJwt.MapClaims{
			"iss": APIServerIssuer,
		}

		// data 实际
		if userInfo, ok := data.(*user.User); ok {
			claims["sub"] = userInfo.Name     // 主题名称，即用户名
			claims["iat"] = time.Now().Unix() // 签发时间
		}

		return claims
	}

}

// 登录回调函数
func authenticator() func(c *gin.Context) (interface{}, error) {

	return func(c *gin.Context) (interface{}, error) {

		var info loginInfo
		var err error

		if c.Request.Header.Get("Authorization") != "" {
			info, err = headBind(c)
			if err != nil {
				return nil, err
			}
		} else {
			info, err = bodyBind(c)
			if err != nil {
				return nil, err
			}
		}
		if err != nil {
			return "", ginJwt.ErrFailedAuthentication
		}

		// 根据 user_name 获取
		userinfo, err := store.GetFactory().User().GetUserByName(c, info.Username)
		if err != nil {
			logrus.Errorf("get user_ information failed: %s", err.Error())

			return "", ginJwt.ErrFailedAuthentication
		}

		// 验证密码
		err = userinfo.Compare(info.Password)
		if err != nil {
			return "", err
		}

		// time.Now()
		nowTime := time.Time{}
		nowTime = time.Now()
		userinfo.LoginedAt = &nowTime

		// 更新登录时间
		_ = store.GetFactory().User().UpdateUser(c, userinfo)
		return userinfo, nil

	}
}

// 通过 head 验证 结构为 Authorization: Basic base64编码的username:password
func headBind(c *gin.Context) (loginInfo, error) {

	basicStr := c.Request.Header.Get("Authorization")

	auths := strings.SplitN(basicStr, " ", 2)
	if len(auths) != 2 || auths[0] != "Basic" {
		return loginInfo{}, ginJwt.ErrFailedAuthentication
	}

	infoStr := base64.StdEncoding.EncodeToString([]byte(auths[1]))
	info := strings.Split(infoStr, ":")
	if len(info) != 2 {
		return loginInfo{}, ginJwt.ErrFailedAuthentication
	}

	return loginInfo{info[0], info[1]}, nil
}

// 通过 body
func bodyBind(c *gin.Context) (loginInfo, error) {

	info := loginInfo{}
	err := c.ShouldBindJSON(&info)
	if err != nil {
		return info, ginJwt.ErrFailedAuthentication
	}

	return info, nil
}

/*

 Header:
      类型声明 : JWT
      声明加密的算法: HMAC SHA256
      密钥Id

 Payload:
      JWT标准注册的声明 (可选)
      公共的声明
      私有的声明

 Signature:
     Header base64. Payload base64 加密后的值
     Salt加密盐值


JWT标准注册的声明：
  iss: JWT token 签发者，其值为大小写敏感的字符串或者Uri
  sub: 主题, sub可用来鉴别一个用户
  exp: JWT 过期时间
  aud: 接收JWT Token 的一方，其值为大小写敏感的字符串或者Uri
  iat: JWT Token 签发时间
  nbf: JWT Token 生效时间
  jti: JWT Token ID ,令牌唯一标识符，通常用于一次性消费的Token
*/
