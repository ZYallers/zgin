package types

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type RestHandlers []RestHandler

type Restful map[string]RestHandlers

type RestHandler struct {
	Sort    int             // Sort
	Signed  bool            // Signed
	Logged  bool            // Logged
	Path    string          // Path
	Version string          // Version
	Http    string          // Http
	Method  string          // Method
	Handler IController     // Handler
	Https   map[string]byte // Https
}

func (rh RestHandlers) Len() int {
	// 返回传入数据的总数
	return len(rh)
}

func (rh RestHandlers) Swap(i, j int) {
	// 两个对象满足Less()则位置对换
	// 表示执行交换数组中下标为i的数据和下标为j的数据
	rh[i], rh[j] = rh[j], rh[i]
}

func (rh RestHandlers) Less(i, j int) bool {
	// 按字段比较大小,此处是降序排序
	// 返回数组中下标为i的数据是否小于下标为j的数据
	return rh[i].Sort > rh[j].Sort
}

func (h *RestHandler) GetHttps() map[string]byte {
	return h.Https
}

func (h *RestHandler) CallMethod(ctx *gin.Context) {
	/*v := reflect.ValueOf(h.Handler)
	ptr := reflect.New(v.Type().Elem())
	ptr.Elem().Set(v.Elem())
	ptr.Interface().(IController).SetContext(ctx)
	//ptr.Elem().FieldByName("Ctx").Set(reflect.ValueOf(ctx))
	ptr.MethodByName(h.Method).Call(nil)*/

	ptr := reflect.New(reflect.ValueOf(h.Handler).Elem().Type())
	ptr.Interface().(IController).SetContext(ctx)
	ptr.MethodByName(h.Method).Call(nil)
}
