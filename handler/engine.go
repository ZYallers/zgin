package handler

import (
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WithNoRoute() option.App {
	return func(app *types.App) {
		app.Server.Http.Handler.(*gin.Engine).NoRoute(func(ctx *gin.Context) {
			ctx.Abort()
			ctx.Header(consts.JsonContentTypeKey, consts.JsonContentTypeValue)
			ctx.Status(http.StatusOK)
			ctx.Writer.WriteString(consts.PageNotFoundContent)
		})
	}
}

func WithRootPath() option.App {
	return func(app *types.App) {
		app.Server.Http.Handler.(*gin.Engine).GET("/", pingHandler)
	}
}

func WithPing() option.App {
	return func(app *types.App) {
		app.Server.Http.Handler.(*gin.Engine).GET("/ping", pingHandler)
	}
}

func WithHealth() option.App {
	return func(app *types.App) {
		app.Server.Http.Handler.(*gin.Engine).GET("/health", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, `ok`)
		})
	}
}

func pingHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, `pong`)
}
