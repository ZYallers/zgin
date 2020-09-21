package zgin

import (
	"flag"
	app "github.com/ZYallers/zgin/application"
	"github.com/ZYallers/zgin/library/logger"
	"github.com/ZYallers/zgin/library/restful"
	"github.com/ZYallers/zgin/library/router"
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
		app.DebugStack = true
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DisableConsoleColor()
	app.Engine = gin.New()
	app.Logger = logger.AppLogger()
}

func SetRouter(restApi *restful.Rest, sessionClient *redis.Client) {
	r := router.NewRouter(app.Engine, app.Logger, app.DebugStack)
	if sessionClient != nil {
		app.Session.Client = sessionClient
	}
	r.RegisterRestApi(restApi).GlobalMiddleware().GlobalHandlerRegister()
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
