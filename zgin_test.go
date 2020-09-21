package zgin

import (
	"github.com/ZYallers/zgin/library/restful"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	EnvInit()
	SetRouter(restful.Api, nil)
	ListenAndServe(5*time.Second, 10*time.Second, 10*time.Second)
}
