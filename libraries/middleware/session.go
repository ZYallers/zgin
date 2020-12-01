package middleware

import (
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"time"
)

const tokenValKeyPrefix = "ci_session:"

func regenSessionData(ctx *gin.Context) {
	defer tool.SafeDefer()

	var client *redis.Client
	if app.Session.Client != nil {
		client = app.Session.Client()
	}
	if client == nil {
		return
	}

	var token string
	if token = queryPostForm(ctx, app.Session.TokenKey); token == "" {
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
			client.Set(tokenValKeyPrefix+token, tool.PhpSerialize(vars),
				time.Duration(app.Session.Expiration)*time.Second)
		}
	}
}

func sessionData(token string) map[string]interface{} {
	var client *redis.Client
	if app.Session.Client != nil {
		client = app.Session.Client()
	}
	if client == nil {
		return nil
	}
	if str, _ := client.Get(tokenValKeyPrefix + token).Result(); str != "" {
		return tool.PhpUnserialize(str)
	}
	return nil
}

func parseSessionToken(ctx *gin.Context) {
	var token string
	if token = queryPostForm(ctx, app.Session.TokenKey); token == "" {
		return
	}
	if vars := sessionData(token); vars != nil {
		ctx.Set(app.Session.DataKey, vars)
	}
}
