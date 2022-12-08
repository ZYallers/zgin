package main

import (
	"fmt"
	"github.com/ZYallers/zgin"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/example/route"
	"github.com/ZYallers/zgin/handler"
	"github.com/ZYallers/zgin/helper/config"
	"github.com/ZYallers/zgin/middleware"
	"github.com/ZYallers/zgin/option"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
	if err := config.ReadFile(); err != nil {
		panic(fmt.Errorf("read config file error: %s", err))
	}
	app := zgin.New(
		option.WithMode(consts.DevMode),
	)
	app.Run(
		handler.WithNoRoute(),
		handler.WithHealth(),
		middleware.WithZapRecovery(),
		middleware.WithZapLogger(),
		handler.WithExpVar(),
		handler.WithPrometheus(),
		handler.WithSwagger(),
		handler.WithPProf(),
		middleware.WithRestCheck(route.Restful),
	)
}
