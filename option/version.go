package option

import (
	"github.com/ZYallers/zgin/types"
)

func WithVersionKey(key string) App {
	return func(app *types.App) {
		app.Version.Key = key
	}
}

func WithVersionLatest(latest string) App {
	return func(app *types.App) {
		app.Version.Latest = latest
	}
}
