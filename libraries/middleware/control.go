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

func AuthCheck(api *app.Restful) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var handler *app.RtHd
		if handler = versionCompare(ctx, api); handler == nil {
			//page404Handler(ctx)
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "msg": "page not found"})
			return
		}

		cp := ctx.Copy()
		errorChan, doneChan, checkNum := make(chan interface{}, 2), make(chan interface{}, 3), 3
		defer func() {
			tool.SafeCloseChan(errorChan)
			tool.SafeCloseChan(doneChan)
		}()

		// 签名验证
		go func(ctx *gin.Context) {
			defer tool.SafeDefer()
			if handler.Signed && !signCheck(ctx) {
				tool.SafeSendChan(errorChan, gin.H{"code": http.StatusForbidden, "msg": "signature error"})
			} else {
				tool.SafeSendChan(doneChan, 1)
			}
		}(cp)

		// 登录验证
		go func(ctx *gin.Context) {
			defer tool.SafeDefer()
			if handler.Logged && !loginCheck(ctx) {
				tool.SafeSendChan(errorChan, gin.H{"code": http.StatusUnauthorized, "msg": "please log in and operate again"})
			} else {
				tool.SafeSendChan(doneChan, 1)
			}
		}(cp)

		// token解析
		go func(ctx *gin.Context) {
			defer tool.SafeDefer()
			if !handler.Logged && handler.ParAck {
				parseSessionToken(ctx)
			} else {
				tool.SafeSendChan(doneChan, 1)
			}
		}(cp)

		var errObj interface{}
	LOOP:
		for {
			select {
			case <-time.After(app.HttpServerWaitTimeout):
				errObj = gin.H{"code": http.StatusGatewayTimeout, "msg": "server processing request timed out"}
				break LOOP
			case resp, ok := <-errorChan:
				if ok && resp != nil {
					errObj = resp
					break LOOP
				}
			case resp, ok := <-doneChan:
				if ok && resp.(int) == 1 {
					if checkNum--; checkNum == 0 {
						break LOOP
					}
				}
			}
		}

		if errObj != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, errObj)
			return
		}

		handler.Handler(ctx)
		ctx.Abort()
	}
}

// versionCompare
func versionCompare(ctx *gin.Context, api *app.Restful) *app.RtHd {
	var (
		exist    bool
		handlers []app.RtHd
	)
	if handlers, exist = (*api)[ctx.Request.URL.Path[1:]]; !exist {
		return nil
	}
	version, method := queryPostForm(ctx, `app_version`, app.Version), ctx.Request.Method
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
func loginCheck(ctx *gin.Context) (pass bool) {
	var token string
	if token = queryPostForm(ctx, "sess_token"); token == "" {
		return
	}
	ctx.Set(app.Session.TokenKey, token)
	if vars := sessionData(token); len(vars) > 0 {
		ctx.Set(app.Session.DataKey, vars)
		if userInfo, ok := vars["userinfo"].(map[string]interface{}); ok {
			if userId, ok := userInfo["userid"].(string); ok && userId != "" {
				ctx.Set(app.Session.LoggedUidKey, userId)
				pass = true
			}
		}
	}
	return
}

// page404Handler
func page404Handler(ctx *gin.Context) {
	go func(ctx *gin.Context) {
		tool.SafeDefer()
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
	}(ctx.Copy())
}
