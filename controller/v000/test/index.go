package v000

import (
	"github.com/ZYallers/zgin/libraries/mvcs"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type index struct {
	mvcs.Controller
}

func Index(c *gin.Context) *index {
	i := &index{}
	i.Ctx = c
	return i
}

func (i *index) CheckOk() {
	i.Json(http.StatusOK, "ok", gin.H{
		"mode":        gin.Mode(),
		"system_ip":   tool.SystemIP(),
		"client_ip":   i.Ctx.ClientIP(),
		"request_url": i.Ctx.Request.URL.String(),
	})
}
