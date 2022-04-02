package middleware

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/ZYallers/golib/funcs/php"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"github.com/syyongx/php2go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RestCheck(c types.ICheck, routes types.Restful) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rest *types.RestHandler
		if rest = versionCompare(ctx, c, routes); rest == nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "msg": "page not found"})
			return
		}

		// 签名验证
		if rest.Signed && !signCheck(ctx, c) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusForbidden, "msg": "signature error"})
			return
		}

		// 解析会话
		if vars := parseSession(ctx, c); vars != nil {
			ctx.Set(consts.SessDataKey, vars)
		}

		// 登录验证
		if rest.Logged && !sessionCheck(ctx) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "msg": "please login first"})
			return
		}

		// 调用对应控制器方法
		rest.CallMethod(ctx)
	}
}

func versionCompare(ctx *gin.Context, c types.ICheck, routes types.Restful) *types.RestHandler {
	var handlers []types.RestHandler
	if path := strings.Trim(ctx.Request.URL.Path, "/"); path == "" {
		return nil
	} else {
		if hd, ok := routes[path]; !ok {
			return nil
		} else {
			handlers = hd
		}
	}
	ver, verKey := c.GetVersion()
	version, httpMethod := queryPostForm(ctx, verKey, ver), ctx.Request.Method
	for _, handler := range handlers {
		hmd := handler.GetHttps()
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

func signCheck(ctx *gin.Context, c types.ICheck) bool {
	secretKey, key, timeKey, dev, expiration := c.GetSign()
	sign := queryPostForm(ctx, key)
	if sign == "" {
		return false
	}
	if gin.IsDebugging() && sign == dev {
		return true
	}
	timestampStr := queryPostForm(ctx, timeKey)
	if timestampStr == "" {
		return false
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 0)
	if err != nil {
		return false
	}
	if time.Now().Unix()-timestamp > expiration {
		return false
	}
	hash := md5.New()
	hash.Write([]byte(timestampStr + secretKey))
	md5str := hex.EncodeToString(hash.Sum(nil))
	if sign == base64.StdEncoding.EncodeToString([]byte(md5str)) {
		return true
	}
	return false
}

func parseSession(ctx *gin.Context, c types.ICheck) map[string]interface{} {
	fn, key, prefix, _ := c.GetSession()
	if fn == nil {
		return nil
	}
	client := fn()
	if client == nil {
		return nil
	}
	token := queryPostForm(ctx, key)
	if token == "" {
		return nil
	}
	str, _ := client.Get(prefix + token).Result()
	if str == "" {
		return nil
	}
	return php.Unserialize(str)
}

func sessionCheck(ctx *gin.Context) bool {
	if vars, ok := ctx.Get(consts.SessDataKey); ok && vars != nil {
		return true
	}
	return false
}

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
