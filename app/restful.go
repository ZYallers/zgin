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

func RtFn(cont interface{}, methodName string) gin.HandlerFunc {
	valueOf := reflect.ValueOf(cont)
	parentValueOf := valueOf.Elem().FieldByName("Controller")
	return func(ctx *gin.Context) {
		parentValueOf.FieldByName("Ctx").Set(reflect.ValueOf(ctx))
		valueOf.MethodByName(methodName).Call(nil)
	}
}
