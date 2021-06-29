package middleware

import (
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const keyPrefix = "ci_session:"

func getSessionClient() *redis.Client {
	if app.Session.GetClientFunc != nil {
		return app.Session.GetClientFunc()
	}
	return nil
}

func sessionData(token string) map[string]interface{} {
	client := getSessionClient()
	if client == nil {
		return nil
	}

	if str, _ := client.Get(keyPrefix + token).Result(); str != "" {
		return tool.PhpUnserialize(str)
	}

	return nil
}

func parseSessionToken(ctx *gin.Context, token string) {
	if vars := sessionData(token); vars != nil {
		ctx.Set(app.Session.DataKey, vars)
	}
}
