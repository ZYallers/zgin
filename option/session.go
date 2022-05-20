package option

import (
	"github.com/ZYallers/zgin/types"
	"github.com/go-redis/redis"
)

func WithSessionKey(key string) App {
	return func(app *types.App) {
		app.Session.Key = key
	}
}

func WithSessionKeyPrefix(keyPrefix string) App {
	return func(app *types.App) {
		app.Session.KeyPrefix = keyPrefix
	}
}

func WithSessionExpiration(expiration int64) App {
	return func(app *types.App) {
		app.Session.Expiration = expiration
	}
}

func WithSessionClientFunc(fn func() *redis.Client) App {
	return func(app *types.App) {
		app.Session.ClientFunc = fn
	}
}
