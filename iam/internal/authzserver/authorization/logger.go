package authorization

import "github.com/ory/ladon"

type AuditLogger struct {
	client AuthorizationInterface
}

func NewAuditLogger(client AuthorizationInterface) ladon.AuditLogger {
	return &AuditLogger{
		client,
	}
}

// LogRejectedAccessRequest 记录被拒绝的授权请求的日志
func (al *AuditLogger) LogRejectedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
	al.client.LogRejectedAccessRequest(request, pool, deciders)
	// log.Debug("subject access review rejected", log.Any("request", r), log.Any("deciders", d))
}

// LogGrantedAccessRequest 记录被允许的授权请求的日志
func (al *AuditLogger) LogGrantedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
	al.client.LogGrantedAccessRequest(request, pool, deciders)
}
