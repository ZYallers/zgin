package zgin

import (
	"github.com/ZYallers/zgin/app"
	v000T "github.com/ZYallers/zgin/controller/v000/test"
	"github.com/ZYallers/zgin/libraries/expvar"
	"github.com/ZYallers/zgin/libraries/prometheus"
	"github.com/ZYallers/zgin/libraries/swagger"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
	"time"
)

var Api = &app.Restful{
	"expvar":    {{Method: map[string]byte{http.MethodGet: 1}, Handler: expvar.RunningStatsHandler}},
	"metrics":   {{Method: map[string]byte{http.MethodGet: 1}, Handler: prometheus.ServerHandler}},
	"swag/json": {{Method: map[string]byte{http.MethodGet: 1}, Handler: swagger.DocsHandler}},
	"test/isok": {{Method: map[string]byte{http.MethodGet: 1}, Handler: func(c *gin.Context) { v000T.Index(c).CheckOk() }}},
}

func TestRun(t *testing.T) {
	EnvInit()
	//SessionClientRegister(nil)
	MiddlewareRegister(Api)
	ListenAndServe(10*time.Second, 15*time.Second, app.HttpServerShutDownTimeout)
}
