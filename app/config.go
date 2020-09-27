package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var (
	HttpServerAddr *string
	Engine         *gin.Engine
	Logger         *zap.Logger
	RobotEnable    bool

	Name                            = "zgin"
	Version                         = "0.0.0"
	VersionKey                      = "app_version"
	HttpServerDefaultAddr           = "0.0.0.0:9010"
	LogDir                          = "/apps/logs/go/zgin"
	ErrorRobotToken                 = ""
	GracefulRobotToken              = ""
	TokenKey                        = ""
	SessionTokenKey                 = "sess_token"
	SignTimeExpiration        int64 = 60
	DevModeSign                     = "hxs-gin-dev"
	LogMaxTimeout                   = 3 * time.Second
	SendMaxTimeout                  = 5 * time.Second
	HttpServerShutDownTimeout       = 15 * time.Second
	HttpServerWaitTimeout           = 10 * time.Second
)
