package store

import (
	pb "iam/pkg/proto/apiserver/v1"
)

// SecretStore 定义密钥相关方法
type SecretStore interface {
	List() (map[string]*pb.SecretInfo, error) // 通过 grpc 获取所有的密钥信息
}
