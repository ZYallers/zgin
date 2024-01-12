package middleware

import (
	"net/http"

	"github.com/ZYallers/golib/funcs/php"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
)

func ParseSession(check types.ICheck) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//defer func(t time.Time) { fmt.Println("ParseSession runtime:", time.Now().Sub(t)) }(time.Now())
		if vars := GetSessionData(ctx, check); vars != nil {
			ctx.Set(consts.SessDataKey, vars)
		}
	}
}

func GetSessionData(ctx *gin.Context, check types.ICheck) map[string]interface{} {
	fn, key, prefix, _ := check.GetSession()
	if fn == nil {
		return nil
	}
	token := QueryPostForm(ctx, key)
	if token == "" {
		return nil
	}
	client := fn()
	if client == nil {
		return nil
	}
	val, _ := client.Get(prefix + token).Result()
	if val == "" {
		return nil
	}
	return php.Unserialize(val)
}

func LoginCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//defer func(t time.Time) { fmt.Println("LoginCheck runtime:", time.Now().Sub(t)) }(time.Now())
		if handler := GetRestHandler(ctx); handler == nil {
			AbortWithJson(ctx, http.StatusUnauthorized, "login handler not found")
			return
		} else {
			if handler.Login {
				if !verifyLogin(ctx) {
					AbortWithJson(ctx, http.StatusUnauthorized, "please login first")
				}
			}
		}
	}
}

func verifyLogin(ctx *gin.Context) bool {
	if vars, ok := ctx.Get(consts.SessDataKey); ok && vars != nil {
		return true
	}
	return false
}
