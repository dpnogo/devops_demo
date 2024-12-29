package middleware

import "github.com/gin-gonic/gin"

/*

  login ------>   api网关验证身份 / 授权 等  ----->  后端server

*/

// iam-apiserver 创建的密钥对是提供给 iam-authorization-server 使用的
// 调用 iam-authorization-server 提供的 RESTful API 接口：/v1/authorization，来进行资源授权
// API 调用比较适合采用的认证方式是 Bearer 认证

// 注意： /v1/authz也可以直接注册到 API 网关中。在实际的 Go 项目开发中，推荐的一种方式。为了展示实现 Bearer 认证的过程，iam-authorization-server 也实现了 Bearer 认证。

// 因为数据库的查询操作延时高，会导致 API 接口延时较高，所以不太适合用在数据流组件中。另外一种是将密码或密钥缓存在内存中，这样请求到来时，就可以直接从内存中查询，从而提升查询速度，提高接口性能
// 将密码或密钥缓存在内存中时，就要考虑内存和数据库的数据一致性，这会增加代码实现的复杂度。
// 因为管控流组件对性能延时要求不那么敏感(iam-apiserver)，而数据流组件则一定要实现非常高的接口性能，

// AuthStrategy 身份验证策略
// auto 策略 根据Head头的 Authorization: Basic xxx 和 Authorization: Bearer xxx 自动选择Basic 认证还是 Bearer 认证
// basic 策略
// jwt 策略 --> bearer 策略
// cache 策略 该策略其实是一个 Bearer 认证的实现，Token 采用了 JWT 格式，因为 Token 中的密钥 ID 是从内存中获取的，所以叫 Cache 认证。
type AuthStrategy interface {
	Auth() gin.HandlerFunc
}

// AuthPolicy 认证策略
type AuthPolicy struct {
	AuthStrategy AuthStrategy
}

// SetPolicy 设置策略
func (ap *AuthPolicy) SetPolicy(as AuthStrategy) {
	ap.AuthStrategy = as
}

// AuthFunc 使用上面设置的策略认证
func (ap *AuthPolicy) AuthFunc() gin.HandlerFunc {
	return ap.AuthStrategy.Auth()
}
