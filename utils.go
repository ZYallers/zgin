package zgin

import (
	"fmt"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/handler"
	"github.com/ZYallers/zgin/helper/dingtalk"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"time"
)

var DefaultApp = &App{
	Name: "demo",
	Mode: consts.DevMode,
	Version: &Version{
		Latest: "1.0.0",
		Key:    "app_version",
	},
	Session: &Session{
		Key:        "sess_token",
		KeyPrefix:  "ci_session:",
		Expiration: 6 * 30 * 86400,
	},
	Logger: &Logger{
		Dir:         "apps/logs/go/demo",
		LogTimeout:  3 * time.Second,
		SendTimeout: 5 * time.Second,
	},
	Sign: &Sign{
		SecretKey:  "1234!@#$",
		Key:        "sess_token",
		TimeKey:    "utime",
		Dev:        "zgin-dev-sign",
		Expiration: 60,
	},
	Server: &Server{
		Addr:            "0.0.0.0:9999",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    15 * time.Second,
		ShutDownTimeout: 15 * time.Second,
	},
}

func LoadJsonFile(args ...string) {
	relativePath, configName, configType := ".", "config", "json"
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

func New() *App {
	app := *DefaultApp
	for key, value := range viper.GetStringMapString("app") {
		switch key {
		case "name":
			app.Name = value
		case "log_dir":
			app.Logger.Dir = value
		case "http_addr":
			app.Server.Addr = value
		case "version":
			app.Version.Latest = value
		case "version_key":
			app.Version.Key = value
		case "sign_key":
			app.Sign.Key = value
		case "sign_time_key":
			app.Sign.TimeKey = value
		case "sign_expiration":
			i, _ := strconv.Atoi(value)
			app.Sign.Expiration = int64(i)
		case "session_key":
			app.Session.Key = value
		case "session_key_prefix":
			app.Session.KeyPrefix = value
		case "session_expiration":
			i, _ := strconv.Atoi(value)
			app.Session.Expiration = int64(i)
		}
	}
	app.Server.Http = &http.Server{
		Handler:      gin.New(),
		Addr:         app.Server.Addr,
		ReadTimeout:  app.Server.ReadTimeout,
		WriteTimeout: app.Server.WriteTimeout,
	}
	logger.SetLoggerDir(app.Logger.Dir)
	return &app
}

func PProfHandler(eg *gin.Engine) {
	handler.PProfRegister(eg)
}

func ExpVarHandler(eg *gin.Engine) {
	eg.GET("/expvar", handler.ExpHandler)
}

func PrometheusHandler(eg *gin.Engine, name string) {
	eg.GET("/metrics", func(ctx *gin.Context) {
		handler.PrometheusHandler(ctx, name)
	})
}

func SwagHandler(eg *gin.Engine) {
	eg.GET("/swag/json", handler.SwagHandler)
}

func StatsVizHandler(eg *gin.Engine, accounts gin.Accounts) {
	relativePath := "/statsviz"
	if accounts == nil {
		eg.GET(relativePath+"/*filepath", handler.StatsHandler)
		return
	}
	eg.Group(relativePath, gin.BasicAuth(accounts)).GET("/*filepath", handler.StatsHandler)
}
