package option

import (
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
)

type App func(app *types.App)

func WithName(name string) App {
	return func(app *types.App) {
		app.Name = name
	}
}

func WithMode(mode string) App {
	return func(app *types.App) {
		app.Mode = mode
		if app.Mode == consts.DevMode {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
	}
}

func WithVersion(ver *types.Version) App {
	return func(app *types.App) {
		app.Version = ver
	}
}

func WithLogger(logger *types.Logger) App {
	return func(app *types.App) {
		app.Logger = logger
	}
}

func WithServer(s *types.Server) App {
	return func(app *types.App) {
		app.Server = s
	}
}

func WithSign(s *types.Sign) App {
	return func(app *types.App) {
		app.Sign = s
	}
}

func WithSession(s *types.Session) App {
	return func(app *types.App) {
		app.Session = s
	}
}
