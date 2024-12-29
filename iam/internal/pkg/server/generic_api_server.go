package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"iam/internal/pkg/middleware"
	appCli "iam/pkg/app/cli"
	"iam/pkg/core"
	"log"
	"net/http"
	"strings"
	"time"
)

// 通过 config --> 进行构建通用服务

// GenericAPIServer 通用API服务器
type GenericAPIServer struct {
	InsecureServing *InsecureServingInfo // http
	SecureServing   *SecureServing       // https
	ShutdownTimeout time.Duration        // Shutdown 等待 X second 进行退出
	middlewares     []string

	*gin.Engine
	healthz         bool // 是否添加检查
	enableProfiling bool // 开启分析 即 pprof
	enableMetrics   bool // 开启 metrics
	// 后续根据需要进行添加功能
	insecureServer, secureServer *http.Server
}

// 初始化该服务应用
func initGenericAPIServer(s *GenericAPIServer) {
	s.Setup()              // debug 时候打印值
	s.InstallMiddlewares() // 加载中间件
	s.InstallAPI()         // 根据配置选项加载所需api
}

// InstallAPI 即health，pprof，metrics根据配置进行启动
func (s *GenericAPIServer) InstallAPI() {

	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}

	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	// 整成中间件形式
	if s.enableMetrics {
		// ginprometheus "github.com/zsais/go-gin-prometheus"
		//prometheus := ginprometheus.NewPrometheus("gin")
		//prometheus.Use(s.Engine)
	}

	s.GET("/version", func(c *gin.Context) {
		core.WriteResponse(c, http.StatusOK, nil, appCli.Get()) // 获取版本信息
	})

}

// Setup gin Mode 设置为 debug 打印值
func (s *GenericAPIServer) Setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// InstallMiddlewares 构建应用中间件
func (s *GenericAPIServer) InstallMiddlewares() {

	s.Use(middleware.RequestId()) // 定义请求request
	// s.Use(middleware.Context())

	// 根据需要进行使用中间件
	for _, v := range s.middlewares {
		// 实现的中间件，才让其进行选择使用
		if _, ok := middleware.Middlewares[v]; ok {
			s.Use(middleware.Middlewares[v])
		}
	}

}

func (s *GenericAPIServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", s.InsecureServing.Address)
	if strings.Contains(s.InsecureServing.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.InsecureServing.Address, ":")[1])
	}

	for {

		// 尝试请求
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		// 测试请求
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Printf("The router has been deployed successfully.\n")

			resp.Body.Close()

			return nil
		}
		time.Sleep(time.Second)

		select {
		case <-ctx.Done(): // 说明上下文过期
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}

}

// Run 构建服务器运行方法
func (s *GenericAPIServer) Run() error {

	s.insecureServer = &http.Server{
		Addr:    s.InsecureServing.Address,
		Handler: s.Engine,
	}

	s.secureServer = &http.Server{
		Addr:    s.SecureServing.Address,
		Handler: s.Engine,
	}

	var eg errgroup.Group

	eg.Go(func() error {
		log.Printf("Start insecure serving addr:%s\n", s.InsecureServing.Address)

		if err := s.insecureServer.ListenAndServe(); err != nil {
			return err
		}

		log.Printf("Stop insecure serving on %s\n", s.InsecureServing.Address)
		return nil
	})

	eg.Go(func() error {
		log.Printf("Start secure serving addr:%s\n", s.SecureServing.Address)

		if err := s.secureServer.ListenAndServeTLS(s.SecureServing.CertKey.CertFile, s.SecureServing.CertKey.KeyFile); err != nil {
			return err
		}

		log.Printf("Stop secure serving on %s\n", s.SecureServing.Address)
		return nil
	})

	// 启动成功后，执行ping保证服务能够正常工作
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	if s.healthz {
		err := s.ping(ctx)
		if err != nil {
			log.Printf("ping the started server err:%v", err)
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

// Close 关闭服务
func (s *GenericAPIServer) Close() {
	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.secureServer.Shutdown(ctx); err != nil {
		log.Printf("Shutdown secure server failed: %s\n", err.Error())
	}

	if err := s.insecureServer.Shutdown(ctx); err != nil {
		log.Printf("Shutdown insecure server failed: %s\n", err.Error())
	}
}
