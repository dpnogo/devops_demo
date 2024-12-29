package authorization

import (
	"github.com/ory/ladon"
)

// 授权接口

type AuthorizationInterface interface {
	Create(*ladon.DefaultPolicy) error // 创建权限
	Update(*ladon.DefaultPolicy) error
	Delete(id string) error
	BatchDelete(ids []string) error
	Get(id string) (*ladon.DefaultPolicy, error)
	List(username string) ([]*ladon.DefaultPolicy, error) // 全部权限

	LogRejectedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) // 授权日志相关
	LogGrantedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies)
}

// AuditLogger: &ladon.AuditLoggerInfo{} 记录授权记录

// AuditLogger 为一个接口 {}

// 将授权记录保存到 redis/db中
// LogRejectedAccessRequest(request *Request, pool Policies, deciders Policies)
// LogGrantedAccessRequest(request *Request, pool Policies, deciders Policies)

// Ladon 支持跟踪一些授权指标，比如 deny、allow、not match、error。

// 跟踪操作
// ladon.Metric 接口  --> 支持跟踪一些授权指标，比如 deny、allow、not match、error

/*
$ curl -s -XPOST -H"Content-Type: application/json" -H"Authorization: Bearer $token" -d'{"metadata":{"name":"authztest"},
 "policy":{
      "description":"One policy to rule them all.",
      "subjects":["users:<peter|ken>","users:maria","groups:admins"],
      "actions":["delete","<create|update>"],
      "effect":"allow",
      "resources":["resources:articles:<.*>","resources:printer"],
      "conditions":{"remoteIP":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'

    http://127.0.0.1:8080/v1/policies


    Subject  string  主体：用户，角色或者服务标识符
    Action   string  操作： 对资源的操作 read write 登
    Resource string  目标资源: 访问/操作的对象  --> 对象
    Context  Context 上下文


*/
