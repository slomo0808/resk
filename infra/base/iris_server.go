package base

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	irisrecover "github.com/kataras/iris/middleware/recover"
	"github.com/sirupsen/logrus"
	"imooc.com/resk/infra"
	"time"
)

var irisApplication *iris.Application

func Iris() *iris.Application {
	return irisApplication
}

type IrisServerStarter struct {
	infra.BaseStarter
}

func (i *IrisServerStarter) Init(ctx infra.StarterContext) {
	// 创建iris application实例
	irisApplication = initIris()
	// 日志组件的配置和扩展
	irisLogger := irisApplication.Logger()
	// 日志组件与logrus 统一
	irisLogger.Install(logrus.StandardLogger())
}

func (i *IrisServerStarter) Start(ctx infra.StarterContext) {
	// 把路由信息打印到控制台
	routers := Iris().GetRoutes()
	for _, r := range routers {
		logrus.Info(r.Trace())
	}
	// 启动iris
	port := ctx.Props().Section("app").Key("server.port").MustString("18080")
	Iris().Run(iris.Addr(":" + port))
}

func (i *IrisServerStarter) StartBlocking() bool {
	return true
}

func initIris() *iris.Application {
	app := iris.New()
	//主要中间件的配置：recover， 日志输出中间件的自定义
	// recover 中间件
	app.Use(irisrecover.New())

	// 日志中间件config
	cfg := logger.Config{
		Status:             true,
		IP:                 true,
		Method:             true,
		Path:               true,
		Query:              true,
		Columns:            false,
		MessageContextKeys: nil,
		MessageHeaderKeys:  nil,
		LogFunc: func(now time.Time, latency time.Duration,
			status, ip, method, path string,
			message interface{},
			headerMessage interface{}) {
			app.Logger().Infof("| %s | %s | %s | %s | %s | %s | %s | %s |",
				now.Format("2006-01-02 15:04:05.000000"),
				latency.String(), status, ip, message, path, headerMessage, message)
		},
		Skippers: nil,
	}
	// 添加日志中间件
	app.Use(logger.New(cfg))
	return app
}
