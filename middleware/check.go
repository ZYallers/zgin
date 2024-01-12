package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
)

const printRouteHandlerFormat = "[GIN-debug] %-6s %-25s --> %s (%d handlers)\n"

func WithRestCheck(routes types.Restful) option.App {
	return func(app *types.App) {
		engine := app.Server.Http.Handler.(*gin.Engine)
		for path, restHandlers := range routes {
			https, sign, login := make(map[string]byte, 0), false, false
			if len(restHandlers) == 1 {
				handler := restHandlers[0]
				https = handler.Https
				sign = handler.Sign
				login = handler.Login
			} else {
				for i := 0; i < len(restHandlers); i++ {
					handler := restHandlers[i]
					for method := range handler.Https {
						if v, ok := https[method]; !ok {
							https[method] = v
						}
					}
					if handler.Sign {
						sign = true
					}
					if handler.Login {
						login = true
					}
				}
			}
			handlers := []gin.HandlerFunc{VersionCompare(app, restHandlers)}
			if sign {
				handlers = append(handlers, SignCheck(app))
			}
			handlers = append(handlers, ParseSession(app))
			if login {
				handlers = append(handlers, LoginCheck())
			}
			handlers = append(handlers, callRestHandler())
			for method := range https {
				switch strings.ToUpper(method) {
				case http.MethodGet:
					engine.GET(path, handlers...)
					printRouteHandler(method, path, handlers)
				case http.MethodPost:
					engine.POST(path, handlers...)
					printRouteHandler(method, path, handlers)
				}
			}
		}
	}
}

func printRouteHandler(method, path string, handlers []gin.HandlerFunc) {
	if gin.IsDebugging() {
		hns := make([]string, 0)
		for _, hd := range handlers {
			fn := handlerName(hd)
			if fn != "" {
				if fns := strings.Split(fn, "."); len(fns) > 2 {
					fn = fns[len(fns)-2]
				}
				hns = append(hns, fn)
			}
		}
		fmt.Printf(printRouteHandlerFormat, method, path, strings.Join(hns, "->"), len(handlers))
	}
}

func handlerName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func callRestHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//defer func(t time.Time) { fmt.Println("callRestHandler runtime:", time.Now().Sub(t)) }(time.Now())
		if handler := GetRestHandler(ctx); handler != nil && handler.Handler != nil && handler.Method != "" {
			ptr := reflect.New(reflect.ValueOf(handler.Handler).Elem().Type())
			if controller, ok := ptr.Interface().(types.IController); ok && controller != nil {
				controller.SetContext(ctx)
				ptr.MethodByName(handler.Method).Call(nil)
			}
		}
	}
}
