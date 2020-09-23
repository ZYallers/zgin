package middleware

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	app "github.com/ZYallers/zgin/application"
	"github.com/ZYallers/zgin/library/restful"
	"github.com/ZYallers/zgin/library/tool"
	"github.com/gin-gonic/gin"
	"github.com/syyongx/php2go"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Dispatch(ege *gin.Engine, rest *restful.Rest) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if regexp.MustCompile(`^\/v[0-9]{3,6}\/.*$`).MatchString(ctx.Request.URL.Path) {
			var handler *restful.RestHandler
			firstIndex := strings.Index(ctx.Request.URL.Path[1:], `/`) + 2
			router := ctx.Request.URL.Path[firstIndex:]
			if handler, _ = versionCompare(router, ctx, rest); handler == nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "msg": "bad request exception"})
				return
			}
			// 签名验证
			if handler.Signed && !signCheck(ctx) {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "msg": "signature error"})
				return
			}
			// 登录验证
			if handler.Logged {
				if !loginCheck(ctx) {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusUnauthorized, "msg": "please log in and operate again"})
					return
				}
			} else {
				if handler.ParAck {
					ackParsing(ctx)
				}
			}
			ctx.Next()
			go regenSessionData(ctx.Copy())
		} else {
			// 版本验证
			if handler, version := versionCompare(ctx.Request.URL.Path[1:], ctx, rest); handler != nil {
				ctx.Request.URL.Path = `/v` + strings.Replace(version, `.`, ``, 2) + ctx.Request.URL.Path
				ege.HandleContext(ctx)
				ctx.Abort()
			} else {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "msg": "page not found"})
			}
		}
	}
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

// versionCompare
func versionCompare(restful string, ctx *gin.Context, rest *restful.Rest) (*restful.RestHandler, string) {
	if handlers, ok := (*rest)[restful]; ok {
		version := queryPostForm(ctx, "app_version", app.Version)
		for _, handler := range handlers {
			if _, ok := handler.Method[ctx.Request.Method]; ok {
				if ilene := len(handler.Version); handler.Version[ilene-1:] == "+" {
					vs := handler.Version[0 : ilene-1]
					if php2go.VersionCompare(version, vs, ">=") {
						return &handler, vs
					}
				} else {
					if php2go.VersionCompare(version, handler.Version, "=") {
						return &handler, version
					}
				}
			}
		}
	}
	return nil, ""
}

func signCheck(ctx *gin.Context) (pass bool) {
	sign := queryPostForm(ctx, "sign")
	if sign == "" {
		return
	}
	timestampStr := queryPostForm(ctx, "utime")
	if timestampStr == "" {
		return
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 0)
	if err != nil {
		return
	}
	if time.Now().Unix()-timestamp > 3600 {
		return
	}
	sign = strings.Trim(sign, " ")
	h := md5.New()
	h.Write([]byte(strconv.FormatInt(timestamp, 10) + app.TokenKey))
	md5str := hex.EncodeToString(h.Sum(nil))
	input := []byte(md5str)
	if sign == base64.StdEncoding.EncodeToString(input) {
		pass = true
	}
	return
}

func loginCheck(ctx *gin.Context) (pass bool) {
	if sessToken := queryPostForm(ctx, "sess_token"); sessToken != "" {
		ctx.Set(app.Session.TokenKey, sessToken)
		if vars := getSessionData(sessToken); len(vars) > 0 {
			ctx.Set(app.Session.DataKey, vars)
			if userInfo, ok := vars["userinfo"].(map[string]interface{}); ok {
				if userId, ok := userInfo["userid"].(string); ok && userId != "" {
					pass = true
					ctx.Set(app.Session.LoggedUidKey, userId)
				}
			}
		}
	}
	return
}

func ackParsing(ctx *gin.Context) {
	if sessToken := queryPostForm(ctx, "sess_token"); sessToken != "" {
		ctx.Set(app.Session.TokenKey, sessToken)
		if vars := getSessionData(sessToken); len(vars) > 0 {
			ctx.Set(app.Session.DataKey, vars)
			if userInfo, ok := vars["userinfo"].(map[string]interface{}); ok {
				if userId, ok := userInfo["userid"].(string); ok && userId != "" {
					ctx.Set(app.Session.LoggedUidKey, userId)
				}
			}
		}
	}
}

func getSessionData(sessToken string) (vars map[string]interface{}) {
	if app.Session.Client != nil {
		if str, _ := app.Session.Client.Get(`ci_session:` + sessToken).Result(); str != "" {
			vars = tool.PhpUnserialize(str)
		}
	}
	return
}

func regenSessionData(ctx *gin.Context) {
	if app.Session.Client == nil {
		return
	}
	var (
		sessToken string
		vars      map[string]interface{}
	)

	if value, ok := ctx.Get(app.Session.TokenKey); !ok {
		return
	} else {
		sessToken = value.(string)
	}

	if value, ok := ctx.Get(app.Session.DataKey); !ok {
		return
	} else {
		vars = value.(map[string]interface{})
	}

	nowTime := time.Now()
	if lastRegen, ok := vars["__ci_last_regenerate"].(int); ok {
		if nowTime.After(time.Unix(int64(lastRegen), 0).Add(app.Session.UpdateDuration)) {
			vars["__ci_last_regenerate"] = nowTime.Unix()
			newCiVars := make(map[string]interface{}, 10)
			if ciVars, ok := vars["__ci_vars"].(map[string]interface{}); ok {
				for k := range ciVars {
					newCiVars[k] = nowTime.Unix() + app.Session.Expiration
				}
				vars["__ci_vars"] = newCiVars
			}
			app.Session.Client.Set(`ci_session:`+sessToken, tool.PhpSerialize(vars), time.Duration(app.Session.Expiration)*time.Second)
		}
	}
}
