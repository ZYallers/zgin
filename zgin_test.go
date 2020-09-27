package zgin

import (
	"github.com/ZYallers/zgin/app"
	v000T "github.com/ZYallers/zgin/controller/v000/test"
	"github.com/ZYallers/zgin/libraries/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

var Api = &app.Restful{
	"expvar":    {{Method: map[string]byte{http.MethodGet: 1}, Handler: handlers.ExpHandler}},
	"metrics":   {{Method: map[string]byte{http.MethodGet: 1}, Handler: handlers.PromHandler}},
	"swag/json": {{Method: map[string]byte{http.MethodGet: 1}, Handler: handlers.SwagHandler}},
	"test/isok": {{Method: map[string]byte{http.MethodGet: 1}, Handler: func(c *gin.Context) { v000T.Index(c).CheckOk() }}},
}

func TestServer(t *testing.T) {
	EnvInit()
	//SessionClientRegister(nil)
	MiddlewareRegister(Api)
	ListenAndServe(app.HttpServerReadTimeout, app.HttpServerWriteTimeout, app.HttpServerShutDownTimeout)
}

func TestServerWithPProf(t *testing.T) {
	EnvInit()
	PProfWebRegister()
	//SessionClientRegister(nil)
	MiddlewareRegister(Api)
	ListenAndServe(app.HttpServerReadTimeout, 60*time.Second, app.HttpServerShutDownTimeout)
}