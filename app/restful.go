package app

import (
	"github.com/gin-gonic/gin"
)

type Restful map[string][]RestHandler

type RestHandler struct {
	Version                string
	Method                 map[string]byte
	Handler                gin.HandlerFunc
	Signed, Logged, ParAck bool
}
