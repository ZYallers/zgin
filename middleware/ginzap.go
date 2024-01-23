package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/ZYallers/golib/funcs/conv"
	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/helper/dingtalk"
	"github.com/ZYallers/zgin/helper/safe"
	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WithZapRecovery() option.App {
	return func(app *types.App) {
		app.Server.Http.Handler.(*gin.Engine).Use(RecoveryWithZap())
	}
}

// RecoveryWithZap returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
func RecoveryWithZap() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// defer func(t time.Time) { fmt.Println("RecoveryWithZap runtime:", time.Now().Sub(t)) }(time.Now())

		// SetPanicOnFault 控制程序在意外（非零）地址出现故障时运行时的行为
		// 此类故障通常是由运行时内存损坏等错误引起的，因此默认响应是使程序崩溃
		// 在不太严重的情况下，使用内存映射文件或不安全的内存操作的程序可能会导致非空地址的错误
		// 允许此类程序请求运行时仅触发恐慌，而不是崩溃
		// 仅适用于当前 goroutine。它返回之前的设置。
		defer debug.SetPanicOnFault(debug.SetPanicOnFault(true))

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

				errs := conv.ToString(err)
				msg := "recovery from panic: " + errs
				reqStr := ctx.GetString(consts.ReqStrKey)
				stacks := string(debug.Stack())

				dingtalk.PushContextMessage(ctx, msg, reqStr, stacks, true)
				logger.Use("recover").Error(msg, zap.String("request", reqStr), zap.String("stack", stacks))

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					_ = ctx.Error(err.(error)) // nolint: errorcheck
					ctx.Abort()
					return
				}

				AbortWithJson(ctx, http.StatusInternalServerError, "server internal error",
					gin.H{"error": errs, "request": reqStr, "stack": strings.Split(stacks, "\n")})
			}
		}()

		reqBytes, _ := httputil.DumpRequest(ctx.Request, true)
		ctx.Set(consts.ReqStrKey, string(reqBytes))
		ctx.Next()
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

func logSend(ctx *gin.Context, runtime, logMaxTime, sendMaxTime time.Duration) {
	defer safe.Defer()
	if len(ctx.Errors) > 0 {
		reqStr := ctx.GetString(consts.ReqStrKey)
		for i, err := range ctx.Errors.Errors() {
			if ctx.Errors[i].Type != gin.ErrorTypeBind {
				logger.Use("context").Error(err)
				dingtalk.PushContextMessage(ctx, err, reqStr, "", true)
			}
		}
	}
	if logMaxTime > 0 && runtime >= logMaxTime {
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
	if sendMaxTime > 0 && runtime >= sendMaxTime {
		msg := fmt.Sprintf("%s take %s to response, exceeding the maximum %s limit", ctx.Request.URL.Path, runtime, sendMaxTime)
		dingtalk.PushContextMessage(ctx, msg, ctx.GetString(consts.ReqStrKey), "", false)
	}
}
