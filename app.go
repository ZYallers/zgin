package zgin

import (
	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/helper/dingtalk"
	"github.com/ZYallers/zgin/helper/grace"
	"github.com/ZYallers/zgin/middleware"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"net/http"
	"strings"
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

type Server struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutDownTimeout time.Duration
	Http            *http.Server
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

func (s *Server) HealthHandler() *Server {
	s.Http.Handler.(*gin.Engine).GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, `ok`)
	})
	return s
}

func (s *Server) NoRouteHandler() *Server {
	s.Http.Handler.(*gin.Engine).NoRoute(func(ctx *gin.Context) {
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
	return s
}

func (a *App) RegisterGlobalMiddleware() *App {
	a.Server.Http.Handler.(*gin.Engine).Use(middleware.RecoveryWithZap(), middleware.LoggerWithZap(a.Logger.LogTimeout, a.Logger.SendTimeout))
	return a
}

func (a *App) RegisterCheckMiddleware(routes types.Restful) *App {
	a.Server.Http.Handler.(*gin.Engine).Use(middleware.RestCheck(a, routes))
	return a
}

func (s *Server) Start() {
	grace.Graceful(s.Http, s.ShutDownTimeout)
}
