package middleware

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/route"
	"github.com/gin-gonic/gin"
	"github.com/syyongx/php2go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// AuthCheck
func AuthCheck(api *route.Restful) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rest *route.RestHandler
		if rest = versionCompare(ctx, api); rest == nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "msg": "page not found"})
			return
		}

		// 签名验证
		if rest.Signed && !signCheck(ctx) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusForbidden, "msg": "signature error"})
			return
		}

		token := queryPostForm(ctx, app.Session.TokenKey)

		// 登录验证
		if rest.Logged && !loginCheck(token) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "msg": "please login first"})
			return
		}

		// 解析会话
		if token != "" {
			parseSessionToken(ctx, token)
		}

		// 调用对应控制器方法
		rest.Handler.SetContext(ctx)
		rest.CallMethod()
	}
}

// versionCompare
func versionCompare(ctx *gin.Context, api *route.Restful) *route.RestHandler {
	var handlers []route.RestHandler

	if path := strings.Trim(ctx.Request.URL.Path, "/"); path == "" {
		return nil
	} else {
		if hd, ok := (*api)[path]; !ok {
			return nil
		} else {
			handlers = hd
		}
	}

	version, httpMethod := queryPostForm(ctx, app.VersionKey, app.Version), ctx.Request.Method
	for _, handler := range handlers {
		hmd := handler.GetHttpMethod()
		if _, exist := hmd[httpMethod]; !exist {
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
