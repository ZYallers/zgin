package gin

import (
	"flag"
	app "github.com/ZYallers/zgin/application"
	"github.com/ZYallers/zgin/library/logger"
	"github.com/ZYallers/zgin/library/restful"
	"github.com/ZYallers/zgin/library/router"
	"github.com/ZYallers/zgin/library/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	app.HttpServerAddr = flag.String("http.addr", app.HttpServerDefaultAddr, "服务监控地址，如：0.0.0.0:9010")
	flag.Parse()

	app.RobotEnable = true
	if os.Getenv("hxsenv") == "development" {
		app.DebugStack = true
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DisableConsoleColor()
	app.Engine = gin.New()
	app.Logger = logger.AppLogger()

	rou := router.NewRouter(app.Engine, app.Logger, app.DebugStack, restful.Api)
	rou.GlobalMiddleware()
	rou.GlobalHandlerRegister()

	srv := &http.Server{
		Addr:         *app.HttpServerAddr,
		Handler:      app.Engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	tool.Graceful(srv, app.Logger, 10*time.Second)
}
