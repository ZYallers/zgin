package zgin

import (
	app "github.com/ZYallers/zgin/application"
	"github.com/ZYallers/zgin/library/restful"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	EnvInit()
	SessionClientRegister(nil)
	MiddlewareRegister(restful.Api)
	ListenAndServe(10*time.Second, 15*time.Second, app.HttpServerShutDownTimeout)
}
