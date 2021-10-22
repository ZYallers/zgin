package v000

import (
	"github.com/ZYallers/zgin/libraries/mvcs"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Index struct {
	mvcs.Controller
	tag struct {
		CheckOk func() `path:"test/isok" http:"get,post"`
		One     func() `path:"test/one" http:"get,post" login:"on"`
		Two     func() `path:"test/two" http:"get,post" sort:"1"`
	}
}

func (i *Index) CheckOk() {
	i.Json(http.StatusOK, "ok", gin.H{
		"mode":      gin.Mode(),
		"system_ip": tool.SystemIP(),
		"client_ip": i.Ctx.ClientIP(),
		"request":   strings.Split(i.DumpRequest(), "\r\n"),
	})
}

func (i *Index) One() {
	i.Json(http.StatusOK, "ok", gin.H{"name": "One"})
}

func (i *Index) Two() {
	i.Json(http.StatusOK, "ok", gin.H{"name": "v000.Index.Two"})
}
