package zgin

import (
	"flag"
	app "github.com/ZYallers/zgin/application"
	"github.com/ZYallers/zgin/library/logger"
	"github.com/ZYallers/zgin/library/middleware"
	"github.com/ZYallers/zgin/library/restful"
	"github.com/ZYallers/zgin/library/tool"
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

func MiddlewareRegister(restApi *restful.Rest) {
	md := []gin.HandlerFunc{
		middleware.RecoveryWithZap(app.Logger),
		middleware.LoggerWithZap(app.Logger),
		middleware.AuthCheck(restApi),
	}
	if app.Session.Client != nil {
		md = append(md, middleware.RegenSessionData())
	}
	app.Engine.Use(md...)
}

func ListenAndServe(readTimeout, writeTimeout, shutdownTimeout time.Duration) {
	srv := &http.Server{
		Addr:         *app.HttpServerAddr,
		Handler:      app.Engine,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	tool.Graceful(srv, app.Logger, shutdownTimeout)
}
