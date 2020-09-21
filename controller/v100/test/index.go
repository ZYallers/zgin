package v100

import (
	"github.com/ZYallers/go-frame/gin/library/mvcs"
	"github.com/ZYallers/go-frame/gin/library/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type index struct {
	mvcs.Controller
}

func Index(ctx *gin.Context) *index {
	i := &index{}
	i.Ctx = ctx
	return i
}

func (i *index) CheckOk() {
	i.Json(http.StatusOK, "ok", gin.H{
		"mode":        gin.Mode(),
		"public_ip":   tool.PublicIP(),
		"system_ip":   tool.SystemIP(),
		"client_ip":   tool.ClientIP(i.Ctx.ClientIP()),
		"request_url": i.Ctx.Request.URL.String(),
	})
}
