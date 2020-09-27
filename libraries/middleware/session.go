package middleware

import (
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/gin-gonic/gin"
	"time"
)

func RegenSessionData() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if app.Session.Client == nil {
			return
		}

		var (
			token string
			vars  map[string]interface{}
		)

		if value, ok := ctx.Get(app.Session.TokenKey); !ok {
			return
		} else {
			token = value.(string)
		}

		if value, ok := ctx.Get(app.Session.DataKey); !ok {
			return
		} else {
			vars = value.(map[string]interface{})
		}

		now := time.Now()
		if lastRegen, ok := vars["__ci_last_regenerate"].(int); ok {
			if now.After(time.Unix(int64(lastRegen), 0).Add(app.Session.UpdateDuration)) {
				vars["__ci_last_regenerate"] = now.Unix()
				newCiVars := make(map[string]interface{}, 10)
				if ciVars, ok := vars["__ci_vars"].(map[string]interface{}); ok {
					for k := range ciVars {
						newCiVars[k] = now.Unix() + app.Session.Expiration
					}
					vars["__ci_vars"] = newCiVars
				}
				app.Session.Client.Set(`ci_session:`+token, tool.PhpSerialize(vars),
					time.Duration(app.Session.Expiration)*time.Second)
			}
		}
	}
}

func sessionData(token string) (vars map[string]interface{}) {
	if app.Session.Client == nil {
		return
	}
	if str, _ := app.Session.Client.Get(`ci_session:` + token).Result(); str != "" {
		vars = tool.PhpUnserialize(str)
	}
	return
}

func parseSessionToken(ctx *gin.Context) {
	if token := queryPostForm(ctx, `sess_token`); token != "" {
		ctx.Set(app.Session.TokenKey, token)
		if vars := sessionData(token); len(vars) > 0 {
			ctx.Set(app.Session.DataKey, vars)
			if userInfo, ok := vars["userinfo"].(map[string]interface{}); ok {
				if userId, ok := userInfo["userid"].(string); ok && userId != "" {
					ctx.Set(app.Session.LoggedUidKey, userId)
				}
			}
		}
	}
}
