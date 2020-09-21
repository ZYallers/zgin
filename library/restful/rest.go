package restful

import "github.com/gin-gonic/gin"

type Rest map[string][]RestHandler

type RestHandler struct {
	Version                string
	Method                 map[string]byte
	Handler                gin.HandlerFunc
	Signed, Logged, ParAck bool
}
