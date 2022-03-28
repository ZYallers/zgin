package zgin

import (
	"context"
	"fmt"
	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/funcs/dingtalk"
	"github.com/ZYallers/zgin/handlers"
	"github.com/ZYallers/zgin/middleware"
	"github.com/ZYallers/zgin/utils/route"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func LoadConfig(args ...string) {
	relativePath, configName, configType := ".", "app", "json"
	argsLen := len(args)
	if argsLen > 0 {
		relativePath = args[0]
	}
	if argsLen > 1 {
		configName = args[1]
	}
	if argsLen > 2 {
		configType = args[2]
	}
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	_, filePath, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filePath), relativePath)
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	if err := viper.ReadInConfig(); err != nil {
		s := fmt.Sprintf("read config file error: %s", err)
		dingtalk.PushSimpleMessage(s, true)
		panic(s)
	}
}

func Initialize() *types.Zgin {
	app := *types.DefaultZgin
	for key, value := range viper.GetStringMap("app") {
		switch key {
		case "name":
			app.Name = value.(string)
		case "version":
			app.Version = value.(string)
		case "versionkey":
			app.VersionKey = value.(string)
		case "httpserveraddr":
			app.HttpServerAddr = value.(string)
		case "tokenkey":
			app.TokenKey = value.(string)
		case "devsign":
			app.DevSign = value.(string)
		case "logdir":
			app.LogDir = value.(string)
		case "signtimeexpiration":
			app.SignTimeExpiration = int64(value.(float64))
		case "logmaxtimeout":
			app.LogMaxTimeout = time.Duration(int64(value.(float64))) * time.Second
		case "sendmaxtimeout":
			app.SendMaxTimeout = time.Duration(int64(value.(float64))) * time.Second
		case "httpserverreadtimeout":
			app.HttpServerReadTimeout = time.Duration(int64(value.(float64))) * time.Second
		case "httpserverwritetimeout":
			app.HttpServerWriteTimeout = time.Duration(int64(value.(float64))) * time.Second
		case "httpservershutdowntimeout":
			app.HttpServerShutDownTimeout = time.Duration(int64(value.(float64))) * time.Second
		case "errorrobottoken":
			app.ErrorRobotToken = value.(string)
		case "gracefulrobottoken":
			app.GracefulRobotToken = value.(string)
		case "session":
			for k, v := range value.(map[string]interface{}) {
				switch k {
				case "tokenkey":
					app.Session.TokenKey = v.(string)
				case "keyprefix":
					app.Session.KeyPrefix = v.(string)
				case "expiration":
					app.Session.Expiration = int64(v.(float64))
				}
			}
		}
	}

	gin.DisableConsoleColor()
	app.Engine = gin.New()
	logger.SetLoggerDir(app.LogDir)
	return &app
}

func SetMode(app *types.Zgin, mode string) {
	app.Mode = mode
	if app.Mode == consts.DevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func HealthHandler(app *types.Zgin) {
	app.Engine.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, `ok`)
	})
}

func NoRouteHandler(app *types.Zgin) {
	app.Engine.NoRoute(func(ctx *gin.Context) {
		reqStr := ctx.GetString(consts.ReqStrKey)
		p := ctx.Request.URL.Path
		logger.Use("404").Info(p,
			zap.String("proto", ctx.Request.Proto),
			zap.String("method", ctx.Request.Method),
			zap.String("host", ctx.Request.Host),
			zap.String("url", ctx.Request.URL.String()),
			zap.String("query", ctx.Request.URL.RawQuery),
			zap.String("clientIP", nets.ClientIP(ctx.ClientIP())),
			zap.Any("header", ctx.Request.Header),
			zap.String("request", reqStr),
		)
		dingtalk.PushContextMessage(ctx, strings.TrimLeft(p, "/")+" page not found", reqStr, "", false)
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "msg": "page not found"})
	})
}

func MiddlewareGlobalRegister(app *types.Zgin) {
	app.Engine.Use(middleware.RecoveryWithZap(), middleware.LoggerWithZap(app))
}

func MiddlewareCustomRegister(app *types.Zgin, api route.Restful) {
	app.Engine.Use(middleware.AuthCheck(app, api))
}

func SessionClientRegister(app *types.Zgin, fn func() *redis.Client) {
	app.Session.GetClientFunc = fn
}

// ***************************************************** Third-party middleware ************************************* //

func PProfRegister(app *types.Zgin) {
	handlers.PProfRegister(app.Engine)
}

func ExpVarRegister(app *types.Zgin) {
	app.Engine.GET("/expvar", handlers.ExpHandler)
}

func PrometheusRegister(app *types.Zgin) {
	app.Engine.GET("/metrics", func(ctx *gin.Context) {
		handlers.PrometheusHandler(ctx, app.Name)
	})
}

func SwaggerRegister(app *types.Zgin) {
	app.Engine.GET("/swag/json", handlers.SwagHandler)
}

func StatsVizRegister(app *types.Zgin, relativePath string, accounts gin.Accounts) {
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

func Serve(app *types.Zgin) {
	srv := &http.Server{
		Addr:         app.HttpServerAddr,
		Handler:      app.Engine,
		ReadTimeout:  app.HttpServerReadTimeout,
		WriteTimeout: app.HttpServerWriteTimeout,
	}
	graceful(srv, app.HttpServerShutDownTimeout)
}

func graceful(srv *http.Server, timeout time.Duration) {
	logAndPushMsg := func(msg string) {
		logger.Use("graceful").Info(msg)
		dingtalk.PushSimpleMessage(msg, true)
	}

	go func() {
		logAndPushMsg(fmt.Sprintf("server(%d) is ready to listen and serve", syscall.Getpid()))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logAndPushMsg(fmt.Sprintf("server listen and serve error: %v", err))
			os.Exit(1)
		}
	}()

	quitChan := make(chan os.Signal, 1)
	// SIGTERM 结束程序(kill pid)(可以被捕获、阻塞或忽略)
	// SIGHUP 终端控制进程结束(终端连接断开)
	// SIGINT 用户发送INTR字符(Ctrl+C)触发
	// SIGQUIT 用户发送QUIT字符(Ctrl+/)触发
	signal.Notify(quitChan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	sign := <-quitChan

	// 保证quitChan将不再接收信号
	signal.Stop(quitChan)

	// 控制是否启用HTTP保持活动，默认情况下始终启用保持活动，只有资源受限的环境或服务器在关闭过程中才应禁用它们
	srv.SetKeepAlivesEnabled(false)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pid := syscall.Getpid()
	logAndPushMsg(fmt.Sprintf("server(%d) is shutting down(%v)...", pid, sign))
	if err := srv.Shutdown(ctx); err != nil {
		logAndPushMsg(fmt.Sprintf("server gracefully shutdown error: %v", err))
	} else {
		logAndPushMsg(fmt.Sprintf("server(%d) has stopped", pid))
	}
}
