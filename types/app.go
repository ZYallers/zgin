package types

import (
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/helper/grace"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"time"
)

var (
	DefaultApp = &App{
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

type Server struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutDownTimeout time.Duration
	Http            *http.Server
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

func (a *App) Run(options ...func(app *App)) {
	for _, opt := range options {
		opt(a)
	}
	grace.Graceful(a.Server.Http, a.Server.ShutDownTimeout, logger.Use(a.Name))
}
