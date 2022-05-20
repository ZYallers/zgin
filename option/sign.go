package option

import (
	"github.com/ZYallers/zgin/types"
)

func WithSignSecretKey(secretKey string) App {
	return func(app *types.App) {
		app.Sign.SecretKey = secretKey
	}
}

func WithSignKey(key string) App {
	return func(app *types.App) {
		app.Sign.Key = key
	}
}

func WithSignTimeKey(timeKey string) App {
	return func(app *types.App) {
		app.Sign.TimeKey = timeKey
	}
}

func WithSignDev(dev string) App {
	return func(app *types.App) {
		app.Sign.Dev = dev
	}
}

func WithSignExpiration(expiration int64) App {
	return func(app *types.App) {
		app.Sign.Expiration = expiration
	}
}
