package zgin

import (
	"github.com/ZYallers/zgin/app"
	v000 "github.com/ZYallers/zgin/controller/v000/test"
	"github.com/ZYallers/zgin/route"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

var api = route.Restful{
	"test/isok": {{Http: "GET,POST", Method: "CheckOk", Handler: &v000.Index{}}},
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
	MiddlewareCustomRegister(api)
	ListenAndServe(app.HttpServerReadTimeout, app.HttpServerWriteTimeout, app.HttpServerShutDownTimeout)
}
