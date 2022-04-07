package zgin

import (
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/middleware"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"time"
)

type Session struct {
	Key        string
	KeyPrefix  string
	Expiration int64
	ClientFunc func() *redis.Client
}

type Version struct {
	Latest, Key string
}

type Logger struct {
	Dir         string
	LogTimeout  time.Duration
	SendTimeout time.Duration
}

type Sign struct {
	SecretKey  string
	Key        string
	TimeKey    string
	Dev        string
	Expiration int64
}

type App struct {
	Name    string
	Mode    string
	Version *Version
	Logger  *Logger
	Server  *Server
	Sign    *Sign
	Session *Session
}

func (a *App) SetMode(mode string) *App {
	a.Mode = mode
	if a.Mode == consts.DevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	return a
}

func (a *App) GetVersion() (version, key string) {
	return a.Version.Latest, a.Version.Key
}

func (a *App) GetSign() (secretKey, key, timeKey, dev string, expiration int64) {
	return a.Sign.SecretKey, a.Sign.Key, a.Sign.TimeKey, a.Sign.Dev, a.Sign.Expiration
}

func (a *App) GetSession() (clientFunc func() *redis.Client, key, prefix string, expiration int64) {
	return a.Session.ClientFunc, a.Session.Key, a.Session.KeyPrefix, a.Session.Expiration
}

func (a *App) RegisterGlobalMiddleware() *App {
	a.Server.Http.Handler.(*gin.Engine).Use(middleware.RecoveryWithZap(), middleware.LoggerWithZap(a.Logger.LogTimeout, a.Logger.SendTimeout))
	return a
}

func (a *App) RegisterCheckMiddleware(routes types.Restful) *App {
	a.Server.Http.Handler.(*gin.Engine).Use(middleware.RestCheck(a, routes))
	return a
}
