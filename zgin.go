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
	developMode = "development"
)

func EnvInit() {
	app.HttpServerAddr = flag.String("http.addr", app.HttpServerDefaultAddr, "服务监控地址，如：0.0.0.0:9010")
	flag.Parse()

	if os.Getenv("hxsenv") == developMode {
		gin.SetMode(gin.DebugMode)
		app.SignTimeExpiration = 3600 // 测试环境utime有效期延长到1小时
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DisableConsoleColor()
	app.Engine = gin.New()
	app.Logger = logger.AppLogger()

	NoRouteHandler()
	HealthHandler()
}

func HealthHandler() {
	app.Engine.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, `ok`)
	})
}

func NoRouteHandler() {
	app.Engine.NoRoute(func(ctx *gin.Context) {
		go middleware.Push404Handler(ctx.Copy())
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "msg": "page not found"})
	})
}

func MiddlewareGlobalRegister() {
	app.Engine.Use(middleware.RecoveryWithZap(app.Logger), middleware.LoggerWithZap(app.Logger))
}

func MiddlewareCustomRegister(api *app.Restful) {
	app.Engine.Use(middleware.AuthCheck(api))
}

func SessionClientRegister(fn func() *redis.Client) {
	if fn != nil {
		app.Session.GetClientFunc = fn
	}
}

// ***************************************************** Third-party middleware ************************************* //

func PProfRegister() {
	handlers.PProfRegister(app.Engine)
}

func ExpVarRegister() {
	app.Engine.GET("/expvar", handlers.ExpHandler)
}

func PrometheusRegister() {
	app.Engine.GET("/metrics", handlers.PromHandler)
}

func SwaggerRegister() {
	app.Engine.GET("/swag/json", handlers.SwagHandler)
}

func StatsVizRegister(relativePath string, accounts gin.Accounts) {
	if relativePath == "" {
		relativePath = "/statsviz"
	}
	if accounts == nil {
		app.Engine.GET(relativePath+"/*filepath", handlers.StatsHandler)
		return
	}
	app.Engine.Group(relativePath, gin.BasicAuth(accounts)).GET("/*filepath", handlers.StatsHandler)
}

// ***************************************************** Server Listen ********************************************** //

func ListenAndServe(readTimeout, writeTimeout, shutdownTimeout time.Duration) {
	srv := &http.Server{
		Addr:         *app.HttpServerAddr,
		Handler:      app.Engine,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	tool.Graceful(srv, app.Logger, shutdownTimeout)
}
