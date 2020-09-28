package zgin

import (
	"github.com/ZYallers/zgin/app"
	v000T "github.com/ZYallers/zgin/controller/v000/test"
	"github.com/gin-gonic/gin"
	"net/http"
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
	if gin.IsDebugging() {
		SwaggerRegister()
		PProfRegister()
		app.HttpServerWriteTimeout = time.Minute
	}
	SessionClientRegister(nil)
	MiddlewareCustomRegister(Api)
	ListenAndServe(app.HttpServerReadTimeout, app.HttpServerWriteTimeout, app.HttpServerShutDownTimeout)
}
