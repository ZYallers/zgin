package middleware

import (
	"fmt"
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/logger"
	"github.com/ZYallers/zgin/libraries/tool"
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

// LoggerWithZap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
func LoggerWithZap(zl *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		go func(ctx *gin.Context, runtime time.Duration) {
			if len(ctx.Errors) > 0 {
				reqStr := ctx.GetString(reqStrKey)
				for _, err := range ctx.Errors.Errors() {
					zl.Error(err)
					tool.PushContextMessage(ctx, err, reqStr, "", true)
				}
			}
			if runtime >= app.LogMaxTimeout {
				logger.Use("timeout").Info(ctx.Request.URL.Path,
					zap.Duration("runtime", runtime),
					zap.String("proto", ctx.Request.Proto),
					zap.String("method", ctx.Request.Method),
					zap.String("host", ctx.Request.Host),
					zap.String("url", ctx.Request.URL.String()),
					zap.String("query", ctx.Request.URL.RawQuery),
					zap.String("clientIP", tool.ClientIP(ctx.ClientIP())),
					zap.Any("header", ctx.Request.Header),
					zap.String("request", ctx.GetString(reqStrKey)),
				)
			}
			if runtime >= app.SendMaxTimeout {
				msg := fmt.Sprintf("%s take %s to response, exceeding the maximum %s limit",
					ctx.Request.URL.Path, runtime, app.SendMaxTimeout)
				tool.PushContextMessage(ctx, msg, ctx.GetString(reqStrKey), "", false)
			}
		}(ctx.Copy(), time.Now().Sub(start))
	}
}

// RecoveryWithZap returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
func RecoveryWithZap(zl *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
				reqStr := ctx.GetString(reqStrKey)
				stacks := string(debug.Stack())

				go func() {
					tool.PushContextMessage(ctx, errMsg, reqStr, stacks, true)
					zl.Error(errMsg, zap.String("request", reqStr), zap.String("stack", stacks))
				}()

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					_ = ctx.Error(err.(error)) // nolint: errorcheck
					ctx.Abort()
					return
				}

				data := gin.H{"error": err}
				if gin.IsDebugging() {
					data["request"] = reqStr
					data["stack"] = stacks
				}
				ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError,
					"msg": "server internal error", "data": data})
			}
		}()

		reqBytes, _ := httputil.DumpRequest(ctx.Request, true)
		ctx.Set(reqStrKey, string(reqBytes))
		ctx.Next()
	}
}
