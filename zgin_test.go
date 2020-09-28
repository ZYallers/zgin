package zgin

import (
	"github.com/ZYallers/zgin/app"
	v000T "github.com/ZYallers/zgin/controller/v000/test"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

var Api = &app.Restful{
	"test/isok": {{Method: map[string]byte{http.MethodGet: 1}, Handler: func(c *gin.Context) { v000T.Index(c).CheckOk() }}},
}

func TestServer(t *testing.T) {
	EnvInit()
	MiddlewareGlobalRegister()
	ExpVarRegister()
	PrometheusRegister()
	SwaggerRegister()
	SessionClientRegister(nil)
	MiddlewareCustomRegister(Api)
	ListenAndServe(app.HttpServerReadTimeout, app.HttpServerWriteTimeout, app.HttpServerShutDownTimeout)
}

func TestServerWithPProf(t *testing.T) {
	EnvInit()
	MiddlewareGlobalRegister()
	ExpVarRegister()
	PrometheusRegister()
	SwaggerRegister()
	app.HttpServerWriteTimeout = 60 * time.Second
	PProfRegister()
	SessionClientRegister(nil)
	MiddlewareCustomRegister(Api)
	ListenAndServe(app.HttpServerReadTimeout, app.HttpServerWriteTimeout, app.HttpServerShutDownTimeout)
}
