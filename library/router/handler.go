package router

import (
	"github.com/ZYallers/zgin/library/logger"
	"github.com/ZYallers/zgin/library/middleware"
	"github.com/ZYallers/zgin/library/restful"
	"github.com/ZYallers/zgin/library/tool"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"strings"
)

type router struct {
	api        *restful.Rest
	engine     *gin.Engine
	logger     *zap.Logger
	debugStack bool
}

func NewRouter(engine *gin.Engine, logger *zap.Logger, debugStack bool) *router {
	return &router{engine: engine, logger: logger, debugStack: debugStack}
}

func (r *router) RegisterRestApi(rest *restful.Rest) *router {
	r.api = rest
	return r
}

func (r *router) GlobalMiddleware() *router {
	r.engine.Use(
		middleware.RecoveryWithZap(r.logger, r.debugStack),
		middleware.Dispatch(r.engine, r.api),
	)
	return r
}

func (r *router) GlobalHandlerRegister() *router {
	r.noRouteHandlerRegister()
	r.versionHandlerRegister(middleware.LoggerWithZap(r.logger))
	return r
}

// adds handlers for NoRoute. It return a 404 code by default.
func (r *router) noRouteHandlerRegister() {
	r.engine.NoRoute(func(ctx *gin.Context) {
		go func(ctx *gin.Context) {
			reqStr := ctx.GetString(`gin-gonic/gin/reqstr`)
			path := ctx.Request.URL.Path
			logger.Use("404").Info(path,
				zap.String("proto", ctx.Request.Proto),
				zap.String("method", ctx.Request.Method),
				zap.String("host", ctx.Request.Host),
				zap.String("url", ctx.Request.URL.String()),
				zap.String("query", ctx.Request.URL.RawQuery),
				zap.String("clientIP", tool.ClientIP(ctx.ClientIP())),
				zap.Any("header", ctx.Request.Header),
				zap.String("request", reqStr),
			)
			tool.PushContextMessage(ctx, strings.TrimLeft(path, "/")+" page not found", reqStr, "", false)
		}(ctx.Copy())
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "msg": "page not found"})
	})
}

func (r *router) versionHandlerRegister(handlers ...gin.HandlerFunc) {
	for path, restHandlers := range *r.api {
		for _, restHandler := range restHandlers {
			version := restHandler.Version
			if versionLen := len(version); version[versionLen-1:] == "+" {
				version = version[0 : versionLen-1]
			}
			version = strings.Join(strings.Split(version, "."), "")
			for method, _ := range restHandler.Method {
				reflect.ValueOf(r.engine).MethodByName(method).CallSlice([]reflect.Value{
					reflect.ValueOf("/v" + version + "/" + path),
					reflect.ValueOf(append(handlers, restHandler.Handler)),
				})
			}
		}
	}
}
