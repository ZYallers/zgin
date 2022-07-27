package zgin

import (
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/helper/config"
	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

var app *types.App

func New(options ...option.App) *types.App {
	app = &types.App{}
	*app = *types.DefaultApp
	for key, value := range config.AppMap() {
		switch key {
		case "name":
			app.Name = cast.ToString(value)
		case "log_dir":
			app.Logger.Dir = cast.ToString(value)
		case "http_addr":
			app.Server.Addr = cast.ToString(value)
		case "version":
			app.Version.Latest = cast.ToString(value)
		case "version_key":
			app.Version.Key = cast.ToString(value)
		case "sign_key":
			app.Sign.Key = cast.ToString(value)
		case "sign_time_key":
			app.Sign.TimeKey = cast.ToString(value)
		case "sign_expiration":
			app.Sign.Expiration = cast.ToInt64(value)
		case "session_key":
			app.Session.Key = cast.ToString(value)
		case "session_key_prefix":
			app.Session.KeyPrefix = cast.ToString(value)
		case "session_expiration":
			app.Session.Expiration = cast.ToInt64(value)
		}
	}
	for _, opt := range options {
		opt(app)
	}
	app.Server.Http = &http.Server{
		Handler:      gin.New(),
		Addr:         app.Server.Addr,
		ReadTimeout:  app.Server.ReadTimeout,
		WriteTimeout: app.Server.WriteTimeout,
	}
	logger.SetLoggerDir(app.Logger.Dir)
	return app
}

func App() *types.App {
	return app
}
