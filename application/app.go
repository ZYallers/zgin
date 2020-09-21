package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	HttpServerAddr *string
	Engine         *gin.Engine
	Logger         *zap.Logger
	RobotEnable    bool
	DebugStack     bool
)

var (
	Name                  = "gin"
	Version               = "1.0.0"
	HttpServerDefaultAddr = "0.0.0.0:9010"
	LogDir                = "/apps/logs/go/gin"
	ErrorRobotToken       = ""
	GracefulRobotToken    = ""
	TokenKey              = ""
)
