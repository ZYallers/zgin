package option

import (
	"github.com/ZYallers/zgin/types"
	"time"
)

func WithLoggerDir(dir string) App {
	return func(app *types.App) {
		app.Logger.Dir = dir
	}
}

func WithLoggerLogTimeout(t time.Duration) App {
	return func(app *types.App) {
		app.Logger.LogTimeout = t
	}
}

func WithLoggerSendTimeout(t time.Duration) App {
	return func(app *types.App) {
		app.Logger.SendTimeout = t
	}
}
