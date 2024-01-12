package option

import (
	"time"

	"github.com/ZYallers/zgin/types"
)

func WithServerAddr(addr string) App {
	return func(app *types.App) {
		app.Server.Addr = addr
	}
}

func WithServerReadTimeout(t time.Duration) App {
	return func(app *types.App) {
		app.Server.ReadTimeout = t
	}
}

func WithServerWriteTimeout(t time.Duration) App {
	return func(app *types.App) {
		app.Server.WriteTimeout = t
	}
}

func WithServerShutDownTimeout(t time.Duration) App {
	return func(app *types.App) {
		app.Server.ShutDownTimeout = t
	}
}
