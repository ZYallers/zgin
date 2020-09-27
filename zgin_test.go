package zgin

import (
	"github.com/ZYallers/zgin/app"
	v000T "github.com/ZYallers/zgin/controller/v000/test"
	"github.com/ZYallers/zgin/libraries/expvar"
	"github.com/ZYallers/zgin/libraries/prometheus"
	"github.com/ZYallers/zgin/libraries/swagger"
	"net/http"
	"testing"
	"time"
)

const (
	hGet = http.MethodGet
)

var Api = &app.Restful{
	"expvar":    {{Method: app.RtMd{hGet: 1}, Handler: expvar.RunningStatsHandler}},
	"metrics":   {{Method: app.RtMd{hGet: 1}, Handler: prometheus.ServerHandler}},
	"swag/json": {{Method: app.RtMd{hGet: 1}, Handler: swagger.DocsHandler}},
	"test/isok": {{Method: app.RtMd{hGet: 1}, Handler: app.RtFn(&v000T.Index{}, "CheckOk")}},
}

func TestRun(t *testing.T) {
	EnvInit()
	SessionClientRegister(nil)
	MiddlewareRegister(Api)
	ListenAndServe(10*time.Second, 15*time.Second, app.HttpServerShutDownTimeout)
}
