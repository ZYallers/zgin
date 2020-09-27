package app

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type Restful map[string][]RtHd

type RtMd map[string]byte

type RtHd struct {
	Version                string
	Method                 RtMd
	Handler                gin.HandlerFunc
	Signed, Logged, ParAck bool
}

func RtFn(controller interface{}, method string) gin.HandlerFunc {
	valueOf := reflect.ValueOf(controller)
	context := valueOf.Elem().FieldByName("Controller").FieldByName("Ctx")
	return func(ctx *gin.Context) {
		context.Set(reflect.ValueOf(ctx))
		valueOf.MethodByName(method).Call(nil)
	}
}
