package middleware

import (
	"fmt"

	"github.com/ZYallers/golib/funcs/conv"
	"github.com/ZYallers/golib/utils/goid"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/golib/utils/trace"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/helper/config"
	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WithOpenTracing() option.App {
	return func(app *types.App) {
		logDir := traceLogDir()
		excludeRoutes := traceExcludeRoutes()
		traceLogger := logger.NewLogger(fmt.Sprintf("%s/%s.log", logDir, app.Name))
		app.Server.Http.Handler.(*gin.Engine).Use(traceHandlerFunc(excludeRoutes, traceLogger))
	}
}

func traceLogDir() string {
	if v := conv.ToString(config.AppValue("trace_log_dir")); v != "" {
		return v
	}
	return "/apps/logs/go/trace"
}

func traceExcludeRoutes() map[string]bool {
	if value := config.AppValue("trace_exclude_routes"); value != nil {
		if val, ok := value.([]interface{}); ok {
			if n := len(val); n > 0 {
				routes := make(map[string]bool, n)
				for _, s := range val {
					routes[conv.ToString(s)] = true
				}
				return routes
			}
		}
	}
	return nil
}

func traceHandlerFunc(excludeRoutes map[string]bool, traceLogger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if excludeRoutes[ctx.Request.URL.Path[1:]] {
			return
		}

		goId := goid.GetString()
		defer trace.DelTraceId(goId)

		traceId := ctx.GetHeader(trace.IdKey)
		if traceId == "" {
			traceId = trace.NewTraceId()
		}
		trace.SetTraceId(goId, traceId)
		ctx.Set(trace.IdKey, traceId)
		ctx.Header(trace.IdKey, traceId)

		traceLogger.Info("",
			zap.String("trace_id", traceId),
			zap.String("host", ctx.Request.Host),
			zap.String("path", ctx.Request.URL.Path),
			zap.String("client_ip", ctx.ClientIP()),
			zap.String("req_raw", ctx.GetString(consts.ReqStrKey)),
		)

		// Must be added ctx.Next(), otherwise there is an issue with the execution order of the defer func above
		ctx.Next()
	}
}
