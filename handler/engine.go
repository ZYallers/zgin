package handler

import (
	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/helper/dingtalk"
	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func WithNoRoute() option.App {
	return func(app *types.App) {
		app.Server.Http.Handler.(*gin.Engine).NoRoute(func(ctx *gin.Context) {
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
