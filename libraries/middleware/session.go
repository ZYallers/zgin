package middleware

import (
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/gin-gonic/gin"
	"time"
)

func regenSessionData(ctx *gin.Context) {
	defer tool.SafeDefer()

	if app.Session.Client == nil {
		return
	}

	var token string
	if token = queryPostForm(ctx, app.SessionTokenKey); token == "" {
		return
	}

	var vars map[string]interface{}
	if value, ok := ctx.Get(app.Session.DataKey); !ok || value == nil {
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

func sessionData(token string) map[string]interface{} {
	if app.Session.Client == nil {
		return nil
	}
	if str, _ := app.Session.Client.Get(`ci_session:` + token).Result(); str != "" {
		return tool.PhpUnserialize(str)
	}
	return nil
}

func parseSessionToken(ctx *gin.Context) {
	if token := queryPostForm(ctx, app.SessionTokenKey); token != "" {
		if vars := sessionData(token); vars != nil {
			ctx.Set(app.Session.DataKey, vars)
		}
	}
}
