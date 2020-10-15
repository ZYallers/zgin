package app

import (
	"github.com/gin-gonic/gin"
)

type Restful map[string][]RestHandler

type RestHandler struct {
	Version string
	Method  map[string]byte
	Handler gin.HandlerFunc
	Signed  bool // 签名验证
	Logged  bool // 登录验证
}
