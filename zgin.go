package zgin

import (
	"flag"
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/handlers"
	"github.com/ZYallers/zgin/libraries/logger"
	"github.com/ZYallers/zgin/libraries/middleware"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"os"
	"time"
)

const (
	developMode = `development`
)

func EnvInit() {
	app.HttpServerAddr = flag.String("http.addr", app.HttpServerDefaultAddr, "服务监控地址，如：0.0.0.0:9010")
	flag.Parse()

	app.RobotEnable = true
	if os.Getenv("hxsenv") == developMode {
		gin.SetMode(gin.DebugMode)
		app.SignTimeExpiration = 3600 // 测试环境utime有效期延长到1小时
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DisableConsoleColor()
	app.Engine = gin.New()
	app.Logger = logger.AppLogger()
}

func SessionClientRegister(cli *redis.Client) {
	if cli != nil {
		app.Session.Client = cli
	}
}

func PProfRegister() {
	if gin.IsDebugging() {
		handlers.PProfRegister(app.Engine)
	}
}

func ExpVarRegister() {
	app.Engine.GET("/expvar", handlers.ExpHandler)
}

func PrometheusRegister() {
	app.Engine.GET("/metrics", handlers.PromHandler)
}

func SwaggerRegister() {
	if gin.IsDebugging() {
		app.Engine.GET("/swag/json", handlers.SwagHandler)
	}
}

func MiddlewareGlobalRegister() {
	app.Engine.Use(middleware.RecoveryWithZap(app.Logger), middleware.LoggerWithZap(app.Logger))
}

func MiddlewareCustomRegister(api *app.Restful) {
	app.Engine.Use(middleware.AuthCheck(api))
}

func ListenAndServe(readTimeout, writeTimeout, shutdownTimeout time.Duration) {
	tool.Graceful(&http.Server{Addr: *app.HttpServerAddr, Handler: app.Engine, ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout}, app.Logger, shutdownTimeout)
}
