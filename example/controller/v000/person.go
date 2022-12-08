package v000

import (
	"github.com/ZYallers/golib/funcs/conv"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Person struct {
	types.Controller
	// tag 标签说明
	// path：接口路由地址
	// ver：支持的版本号，尾部添加"+"符号，表示大于等于当前版本都支持
	// http：支持的请求方式
	// sign：是否要求签名验证，on为开启，不传或其他值为不开启
	// login：是否要求登录验证，on为开启，不传或其他值为不开启
	// sort：排序号，当存在多个方法path一样的时候才有用，sort越大排越前，在符合ver条件下优化调用
	tag struct {
		List func() `path:"person/list" http:"get,post" sign:"off" login:"off" sort:"1"`
		Info func() `path:"person/info" http:"get,post" sign:"off" login:"off"`
	}
}

func (p *Person) List() {
	lists := gin.H{
		"1": gin.H{"id": 1, "name": "peter", "age": 12},
		"2": gin.H{"id": 2, "name": "jack", "age": 11},
		"3": gin.H{"id": 3, "name": "ace", "age": 13},
	}
	p.Json(http.StatusOK, "success", gin.H{"list": lists})
}

func (p *Person) Info() {
	id := conv.ToInt(p.GetQueryPostForm("id"))
	if id <= 0 {
		p.Json(http.StatusNotImplemented, "illegal parameter")
		return
	}
	var info gin.H
	switch id {
	case 1:
		info = gin.H{"id": 1, "name": "peter", "age": 12}
		p.Json(http.StatusOK, "success", gin.H{"info": info})
	default:
		p.Json(http.StatusNoContent, "person does not exist")
	}
}
