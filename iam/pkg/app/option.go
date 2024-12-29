package app

import (
	pkg "iam/pkg/app/cli"
)

// CliOptions 从命令行/配置等读取参数的选项接口
type CliOptions interface {
	Flags() pkg.NamedFlagSets // 获取所有的flag值
	Validate() []error        // 对相关的配置进行校验 同时验证所有的命令行参数/配置文件是否有效
}

// ConfigurableOptions 抽象了从配置文件中读取参数的配置选项。
type ConfigurableOptions interface {
	// ApplyFlags 从命令行或配置文件解析参数到options实例。
	ApplyFlags() []error
}

// CompleteableOptions 抽象了可以完成的选项， 验证参数是否完成，例如 jwt key是否为空，https tls的key和
type CompleteableOptions interface {
	Complete() error
}

// PrintableOptions 抽象可打印的选项
type PrintableOptions interface {
	String() string
}
