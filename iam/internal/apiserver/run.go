package apiserver

import (
	"iam/internal/apiserver/config"
	"iam/internal/apiserver/options"
	"iam/pkg/app"
	"iam/pkg/logger"
)

func NewApp(basename string) *app.App {

	opts := options.NewOptions()
	// 初始化应用框架，内部将cli+config+env并进行合并，并解析到cfg中
	application := app.NewApp(
		"IAM API Server",
		basename,
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
		app.WithOptions(opts),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {

		initLog(opts)

		cfg := config.NewConfig(opts)
		// 创建 .keep-server 服务 并进行运行
		server, err := newApiServer(cfg)
		if err != nil {
			return err
		}

		return server.PrepareServer().Run()
	}
}

// 根据配置初始化 log
func initLog(opts *options.Options) {
	logger.NewLog(logger.LogCfg{
		LogLevel: opts.Log.Level,
	})
}
