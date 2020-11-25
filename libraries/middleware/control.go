package middleware

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/logger"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/gin-gonic/gin"
	"github.com/syyongx/php2go"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// AuthCheck
func AuthCheck(api *app.Restful) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rh *app.RestHandler
		if rh = versionCompare(ctx, api); rh == nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "msg": "page not found"})
			return
		}

		// 签名验证
		if rh.Signed && !signCheck(ctx) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusForbidden, "msg": "signature error"})
			return
		}

		// 登录验证
		if rh.Logged && !loginCheck(ctx) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "msg": "please login first"})
			return
		}

		// 解析sessionToken
		parseSessionToken(ctx)

		rh.Handler(ctx)
		go regenSessionData(ctx.Copy())
	}
}

// versionCompare
func versionCompare(ctx *gin.Context, api *app.Restful) *app.RestHandler {
	var (
		exist    bool
		handlers []app.RestHandler
	)
	if handlers, exist = (*api)[ctx.Request.URL.Path[1:]]; !exist {
		return nil
	}
	version, method := queryPostForm(ctx, app.VersionKey, app.Version), ctx.Request.Method
	for _, handler := range handlers {
		if _, ok := handler.Method[method]; ok {
			if handler.Version == `` || version == handler.Version {
				return &handler
			}
			if le := len(handler.Version); handler.Version[le-1:] == `+` {
				vs := handler.Version[0 : le-1]
				if version == vs {
					return &handler
				}
				if php2go.VersionCompare(version, vs, `gt`) {
					return &handler
				}
			}
		}
	}
	return nil
}

// signCheck
func signCheck(ctx *gin.Context) bool {
	sign := queryPostForm(ctx, "sign")
	if sign == "" {
		return false
	}
	// 开发测试模式下，用固定sign判断
	if gin.IsDebugging() && sign == app.DevModeSign {
		return true
	}
	timestampStr := queryPostForm(ctx, "utime")
	if timestampStr == "" {
		return false
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 0)
	if err != nil {
		return false
	}
	if time.Now().Unix()-timestamp > app.SignTimeExpiration {
		return false
	}
	hash := md5.New()
	hash.Write([]byte(timestampStr + app.TokenKey))
	md5str := hex.EncodeToString(hash.Sum(nil))
	if sign == base64.StdEncoding.EncodeToString([]byte(md5str)) {
		return true
	}
	return false
}

// loginCheck
func loginCheck(ctx *gin.Context) bool {
	if app.Session.Client == nil {
		return false
	}
	var token string
	if token = queryPostForm(ctx, app.Session.TokenKey); token == "" {
		return false
	}
	if vars := sessionData(token); vars != nil {
		return true
	}
	return false
}

// push404Handler
func push404Handler(ctx *gin.Context) {
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
