package v100

import (
	"github.com/ZYallers/zgin/library/mvcs"
	"github.com/ZYallers/zgin/library/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Index struct {
	mvcs.Controller
}

func (i *Index) CheckOk() {
	i.Json(http.StatusOK, "ok", gin.H{
		"mode":        gin.Mode(),
		"public_ip":   tool.PublicIP(),
		"system_ip":   tool.SystemIP(),
		"client_ip":   tool.ClientIP(i.Ctx.ClientIP()),
		"request_url": i.Ctx.Request.URL.String(),
	})
}
