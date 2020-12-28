package middleware

import (
	"github.com/ZYallers/zgin/libraries/logger"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

const (
	reqStrKey = `gin-gonic/gin/reqstr`
)

// Push404Handler
func Push404Handler(ctx *gin.Context) {
	defer tool.SafeDefer()
	reqStr := ctx.GetString(reqStrKey)
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
}

// queryPostForm
func queryPostForm(ctx *gin.Context, keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	if val, ok := ctx.GetQuery(keys[0]); ok {
		return val
	}
	if val, ok := ctx.GetPostForm(keys[0]); ok {
		return val
	}
	if len(keys) == 2 {
		return keys[1]
	}
	return ""
}
