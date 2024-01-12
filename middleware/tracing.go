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
		traceLogDir := conv.ToString(config.AppValue("trace_log_dir"))
		traceLogger := logger.NewLogger(fmt.Sprintf("%s/%s.log", traceLogDir, app.Name))
		app.Server.Http.Handler.(*gin.Engine).Use(func(ctx *gin.Context) {
			goId := goid.GetString()
			traceId := ctx.GetHeader(trace.IdKey)
			if traceId == "" {
				traceId = trace.GenTraceId()
			}
			trace.SetTraceId(goId, traceId)
			ctx.Set(trace.IdKey, traceId)
			ctx.Header(trace.IdKey, traceId)
			traceLogger.Info(app.Name,
				zap.String("trace_id", traceId),
				zap.String("host", ctx.Request.Host),
				zap.String("path", ctx.Request.URL.Path),
				zap.String("client_ip", ctx.ClientIP()),
				zap.String("req_raw", ctx.GetString(consts.ReqStrKey)),
			)
			ctx.Next()
			trace.DelTraceId(goId)
		})
	}
}
