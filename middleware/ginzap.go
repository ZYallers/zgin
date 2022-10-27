package middleware

import (
	"fmt"
	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/helper/dingtalk"
	"github.com/ZYallers/zgin/helper/safe"
	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const sendMaxTimeFormat = "%s take %s to response, exceeding the maximum %s limit"

func WithZapRecovery() option.App {
	return func(app *types.App) {
		app.Server.Http.Handler.(*gin.Engine).Use(RecoveryWithZap())
	}
}

func WithZapLogger() option.App {
	return func(app *types.App) {
		app.Server.Http.Handler.(*gin.Engine).Use(LoggerWithZap(app.Logger.LogTimeout, app.Logger.SendTimeout))
	}
}

// LoggerWithZap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
func LoggerWithZap(logMaxTime, sendMaxTime time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//defer func(t time.Time) { fmt.Println("LoggerWithZap runtime:", time.Now().Sub(t)) }(time.Now())
		start := time.Now()
		ctx.Next()
		go logSend(ctx.Copy(), time.Now().Sub(start), logMaxTime, sendMaxTime)
	}
}

func logSend(ctx *gin.Context, runtime, logTime, sendTime time.Duration) {
	defer safe.Defer()
	if len(ctx.Errors) > 0 {
		reqStr := ctx.GetString(consts.ReqStrKey)
		for _, err := range ctx.Errors.Errors() {
			logger.Use("context").Error(err)
			dingtalk.PushContextMessage(ctx, err, reqStr, "", true)
		}
	}
	if runtime >= logTime {
		logger.Use("timeout").Info(ctx.Request.URL.Path,
			zap.Duration("runtime", runtime),
			zap.String("proto", ctx.Request.Proto),
			zap.String("method", ctx.Request.Method),
			zap.String("host", ctx.Request.Host),
			zap.String("url", ctx.Request.URL.String()),
			zap.String("query", ctx.Request.URL.RawQuery),
			zap.String("clientIP", nets.ClientIP(ctx.ClientIP())),
			zap.Any("header", ctx.Request.Header),
			zap.String("request", ctx.GetString(consts.ReqStrKey)),
		)
	}
	if runtime >= sendTime {
		msg := fmt.Sprintf(sendMaxTimeFormat, ctx.Request.URL.Path, runtime, sendTime)
		dingtalk.PushContextMessage(ctx, msg, ctx.GetString(consts.ReqStrKey), "", false)
	}
}

// RecoveryWithZap returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
func RecoveryWithZap() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//defer func(t time.Time) { fmt.Println("RecoveryWithZap runtime:", time.Now().Sub(t)) }(time.Now())
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				errMsg := fmt.Sprintf("recovery from panic: %v", err)
				reqStr := ctx.GetString(consts.ReqStrKey)
				stacks := string(debug.Stack())

				dingtalk.PushContextMessage(ctx, errMsg, reqStr, stacks, true)
				logger.Use("recover").Error(errMsg, zap.String("request", reqStr), zap.String("stack", stacks))

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					_ = ctx.Error(err.(error)) // nolint: errorcheck
					ctx.Abort()
					return
				}

				data := gin.H{"error": err}
				if gin.IsDebugging() {
					data["request"] = reqStr
					data["stack"] = strings.Split(stacks, "\n")
				}
				AbortWithJson(ctx, http.StatusInternalServerError, "server internal error", data)
			}
		}()

		reqBytes, _ := httputil.DumpRequest(ctx.Request, true)
		ctx.Set(consts.ReqStrKey, string(reqBytes))
		ctx.Next()
	}
}
