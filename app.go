package zgin

import (
	"net/http"

	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/helper/config"
	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
)

var app *types.App

func New(options ...option.App) *types.App {
	app = &types.App{}
	*app = *types.DefaultApp
	for k, v := range config.AppMap() {
		switch k {
		case "name":
			app.Name = v.(string)
		case "log_dir":
			app.Logger.Dir = v.(string)
		case "http_addr":
			app.Server.Addr = v.(string)
		case "version":
			app.Version.Latest = v.(string)
		case "version_key":
			app.Version.Key = v.(string)
		case "sign_key":
			app.Sign.Key = v.(string)
		case "sign_time_key":
			app.Sign.TimeKey = v.(string)
		case "sign_expiration":
			app.Sign.Expiration = int64(v.(float64))
		case "session_key":
			app.Session.Key = v.(string)
		case "session_key_prefix":
			app.Session.KeyPrefix = v.(string)
		case "session_expiration":
			app.Session.Expiration = int64(v.(float64))
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
