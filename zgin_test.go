package zgin

import (
	"github.com/ZYallers/zgin/app"
	v000 "github.com/ZYallers/zgin/controller/v000"
	v110 "github.com/ZYallers/zgin/controller/v110"
	"github.com/ZYallers/zgin/route"
	"github.com/gin-gonic/gin"
	"log"
	"testing"
	"time"
)

var api route.Restful

func init() {
	/*i000 := route.Restful{
		"test/isok": {{Http: "GET,POST", Method: "CheckOk", Handler: &v000.Index{}}},
		"test/one":  {{Http: "GET,POST", Method: "One", Handler: &v000.Index{}}},
		"test/index/two": {
			{Version: "1.1.0+", Http: "GET,POST", Method: "Two", Handler: &v110.Index{}},
			{Http: "GET,POST", Method: "Two", Handler: &v000.Index{}},
		},
	}
	i110 := route.Restful{
		"test/third": {{Version: "1.1.0+", Http: "GET,POST", Method: "Third", Handler: &v110.Index{}}},
	}
	api = route.Merge(i000, i110)*/
	api = route.Register(nil, &v000.Index{}, &v110.Index{})
}

func TestServer(t *testing.T) {
	EnvInit()
	MiddlewareGlobalRegister()
	ExpVarRegister()
	PrometheusRegister()
	if gin.IsDebugging() {
		SwaggerRegister()
		PProfRegister()
		StatsVizRegister("", gin.Accounts{"test": "123456"})
		app.HttpServerWriteTimeout = time.Minute
	}
	SessionClientRegister(nil)
	for path, restHandlers := range api {
		log.Printf("path: %s\n", path)
		for i, resHandler := range restHandlers {
			log.Printf("handlers.%d: %+v\n", i, resHandler)
		}
	}
	MiddlewareCustomRegister(api)
	ListenAndServe(app.HttpServerReadTimeout, app.HttpServerWriteTimeout, app.HttpServerShutDownTimeout)
}
