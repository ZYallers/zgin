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
}

func (i *Index) CheckOk() {
	i.Json(http.StatusOK, "ok", gin.H{
		"mode":      gin.Mode(),
		"system_ip": tool.SystemIP(),
		"client_ip": i.Ctx.ClientIP(),
		"request":   strings.Split(i.DumpRequest(), "\r\n"),
	})
}
