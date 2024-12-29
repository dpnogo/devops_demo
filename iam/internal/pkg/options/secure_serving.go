package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"iam/internal/pkg/server"
	"net"
	"path"
	"strconv"
)

// SecureServingOptions 不安全服务选项
type SecureServingOptions struct {
	BindAddress string   `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int      `json:"bind-port" mapstructure:"bind-port"`
	Required    bool     // Required 设置为true意味着BindPort不能为零。
	ServerCert  CertConf `json:"tls"          mapstructure:"tls"`
}

type CertKey struct {
	CertFile string `json:"cert-file" mapstructure:"cert-file"` // 证书文件
	KeyFile  string `json:"key-file" mapstructure:"key-file"`   // 密钥key
}

type CertConf struct {
	CertKey CertKey `json:"cert-key" mapstructure:"cert-key"`

	// CertDirectory 如果没有明确设置CertFile/KeyFile, CertDirectory指定一个目录来写入生成的证书  优先级 < CertKey
	CertDirectory string `json:"cert-directory" mapstructure:"cert-directory"`
	// PairName 是将与CertDirectory一起使用的名称，用于生成证书和密钥文件名,
	PairName string `json:"pair-name" mapstructure:"pair-name"`
}

func NewSecureServing() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress: "127.0.0.1",
		BindPort:    443,
		ServerCert: CertConf{
			PairName:      "iam",
			CertDirectory: "/var/run/iam",
		},
	}
}

func (svc *SecureServingOptions) Validate() []error {

	var (
		errs []error
	)

	if svc.Required && svc.BindPort < 1 || svc.BindPort > 65535 {
		errs = append(errs, fmt.Errorf("service port must >= 1 and <= 65535, port:%d", svc.BindPort))
	}

	if svc.BindPort < 0 || svc.BindPort > 65535 {
		errs = append(errs, fmt.Errorf("service port must >= 0 and <= 65535, port:%d", svc.BindPort))
	}

	return errs
}

// ApplyTo 应用options的配置到config
func (svc *SecureServingOptions) ApplyTo(config *server.Config) error {
	config.SecureServing = &server.SecureServing{
		Address: net.JoinHostPort(svc.BindAddress, strconv.Itoa(svc.BindPort)),
		CertKey: server.CertKey{
			CertFile: svc.ServerCert.CertKey.CertFile,
			KeyFile:  svc.ServerCert.CertKey.KeyFile,
		},
	}
	return nil
}

func (svc *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {

	var (
		desc string
	)
	if svc.Required {
		desc += fmt.Sprintf("port must >= %d and <= 65535", 1)
	} else {
		desc += fmt.Sprintf("port must >= %d and <= 65535", 0)
	}

	fs.StringVar(&svc.BindAddress, "address", "", "secure_svc_options address flag") // flag 对应的值
	fs.IntVarP(&svc.BindPort, "port", "p", 0, fmt.Sprintf("secure_svc_options port flag, desc:%s", desc))

	// tls 文件相关， 证书 + key
	fs.StringVar(&svc.ServerCert.CertDirectory, "secure.tls.cert-dir", svc.ServerCert.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file are provided, "+
		"this flag will be ignored.")

	fs.StringVar(&svc.ServerCert.PairName, "secure.tls.pair-name", svc.ServerCert.PairName, ""+
		"The name which will be used with --secure.tls.cert-dir to make a cert and key filenames. "+
		"It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key")

	fs.StringVar(&svc.ServerCert.CertKey.CertFile, "secure.tls.cert-key.cert-file", svc.ServerCert.CertKey.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&svc.ServerCert.CertKey.KeyFile, "secure.tls.cert-key.private-key-file",
		svc.ServerCert.CertKey.KeyFile, ""+
			"File containing the default x509 private key matching --secure.tls.cert-key.cert-file.")

}

// Complete 填充任何未设置的字段，这些字段必须具有有效数据
func (svc *SecureServingOptions) Complete() error {
	if svc == nil || svc.BindPort == 0 {
		return nil
	}

	keyCert := &svc.ServerCert.CertKey
	if len(keyCert.CertFile) != 0 || len(keyCert.KeyFile) != 0 {
		return nil
	}

	if len(svc.ServerCert.CertDirectory) > 0 {
		if len(svc.ServerCert.PairName) == 0 {
			return fmt.Errorf("--secure.tls.pair-name is required if --secure.tls.cert-dir is set")
		}
		keyCert.CertFile = path.Join(svc.ServerCert.CertDirectory, svc.ServerCert.PairName+".crt")
		keyCert.KeyFile = path.Join(svc.ServerCert.CertDirectory, svc.ServerCert.PairName+".key")
	}

	return nil
}
