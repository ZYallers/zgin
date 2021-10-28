package route

import (
	"github.com/ZYallers/zgin/libraries/mvcs"
	"github.com/gin-gonic/gin"
	"reflect"
)

type RestHandlers []RestHandler

type Restful map[string]RestHandlers

type RestHandler struct {
	Sort    int              // Sort
	Signed  bool             // Signed
	Logged  bool             // Logged
	Path    string           // Path
	Version string           // Version
	Http    string           // Http
	Method  string           // Method
	Handler mvcs.IController // Handler
	http    map[string]byte  // http
}

// RestHandlers 实现sort SDK中的Interface接口
// Len
func (rh RestHandlers) Len() int {
	// 返回传入数据的总数
	return len(rh)
}

// Swap
func (rh RestHandlers) Swap(i, j int) {
	// 两个对象满足Less()则位置对换
	// 表示执行交换数组中下标为i的数据和下标为j的数据
	rh[i], rh[j] = rh[j], rh[i]
}

// Less
func (rh RestHandlers) Less(i, j int) bool {
	// 按字段比较大小,此处是降序排序
	// 返回数组中下标为i的数据是否小于下标为j的数据
	return rh[i].Sort > rh[j].Sort
}

// GetHttpMethod
func (rh *RestHandler) GetHttpMethod() map[string]byte {
	return rh.http
}

// CallMethod
func (rh *RestHandler) CallMethod(ctx *gin.Context) {
	ptr := reflect.New(reflect.TypeOf(rh.Handler).Elem())
	ptr.Elem().Set(reflect.ValueOf(rh.Handler).Elem())
	ptr.Elem().FieldByName("Ctx").Set(reflect.ValueOf(ctx))
	ptr.MethodByName(rh.Method).Call(nil)
}
