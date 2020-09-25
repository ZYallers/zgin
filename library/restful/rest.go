package restful

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type Rest map[string][]RestHandler

type Mt map[string]byte

type RestHandler struct {
	Version                string
	Method                 Mt
	Handler                gin.HandlerFunc
	Signed, Logged, ParAck bool
}

func Fn(cont interface{}, methodName string) gin.HandlerFunc {
	valueOf := reflect.ValueOf(cont)
	parentValueOf := valueOf.Elem().FieldByName("Controller")
	return func(ctx *gin.Context) {
		parentValueOf.FieldByName("Ctx").Set(reflect.ValueOf(ctx))
		valueOf.MethodByName(methodName).Call(nil)
	}
}
