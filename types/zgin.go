package types

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Zgin struct {
	Name           string
	Version        string
	VersionKey     string
	LogDir         string
	HttpServerAddr string

	Mode                      string
	TokenKey                  string
	DevSign                   string
	ErrorRobotToken           string
	GracefulRobotToken        string
	SignTimeExpiration        int64
	LogMaxTimeout             time.Duration
	SendMaxTimeout            time.Duration
	HttpServerReadTimeout     time.Duration
	HttpServerWriteTimeout    time.Duration
	HttpServerShutDownTimeout time.Duration

	Engine  *gin.Engine
	Session *Session
}

var DefaultZgin = &Zgin{
	Name:                      "zgin",
	Version:                   "1.0.0",
	VersionKey:                "app_version",
	LogDir:                    "apps/logs/go/zgin",
	SignTimeExpiration:        60,
	LogMaxTimeout:             3 * time.Second,
	SendMaxTimeout:            5 * time.Second,
	HttpServerReadTimeout:     10 * time.Second,
	HttpServerWriteTimeout:    15 * time.Second,
	HttpServerShutDownTimeout: 15 * time.Second,
	Session:                   DefaultSession,
}
