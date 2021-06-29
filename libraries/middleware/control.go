package middleware

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/ZYallers/zgin/app"
	"github.com/gin-gonic/gin"
	"github.com/syyongx/php2go"
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

		token := queryPostForm(ctx, app.Session.TokenKey)

		// 登录验证
		if rh.Logged && !loginCheck(token) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "msg": "please login first"})
			return
		}

		// 解析会话
		if token != "" {
			parseSessionToken(ctx, token)
		}

		if rh.Handler != nil {
			rh.Handler(ctx)
		}
	}
}

// versionCompare
func versionCompare(ctx *gin.Context, api *app.Restful) *app.RestHandler {
	var handlers []app.RestHandler

	if path := strings.Trim(ctx.Request.URL.Path, "/"); path == "" {
		return nil
	} else {
		if hd, ok := (*api)[path]; !ok {
			return nil
		} else {
			handlers = hd
		}
	}

	version, method := queryPostForm(ctx, app.VersionKey, app.Version), ctx.Request.Method
	for _, handler := range handlers {
		if handler.Method == nil {
			return nil
		}
		if _, ok := handler.Method[method]; !ok {
			return nil
		}
		if handler.Version == "" || version == handler.Version {
			return &handler
		}
		if le := len(handler.Version); handler.Version[le-1:] == "+" {
			vs := handler.Version[0 : le-1]
			if version == vs {
				return &handler
			}
			if php2go.VersionCompare(version, vs, ">") {
				return &handler
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
func loginCheck(token string) bool {
	if token == "" {
		return false
	}
	if vars := sessionData(token); vars != nil {
		return true
	}
	return false
}
