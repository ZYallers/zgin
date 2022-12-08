package v100

import (
	"github.com/ZYallers/golib/funcs/conv"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Person struct {
	types.Controller
	tag struct {
		List        func() `path:"person/list" ver:"1.0.0+" http:"get,post" sign:"off" login:"off" sort:"2"`
		PrivateInfo func() `path:"person/private/info" ver:"1.0.0+" http:"get,post" sign:"on" login:"off"`
	}
}

func (p *Person) List() {
	lists := gin.H{
		"1": gin.H{"id": 1, "name": "peter", "age": 12, "birthday": "1990-12-11"},
		"2": gin.H{"id": 2, "name": "jack", "age": 11, "birthday": "1989-09-20"},
		"3": gin.H{"id": 3, "name": "ace", "age": 13, "birthday": "1991-07-15"},
	}
	p.Json(http.StatusOK, "success", gin.H{"list": lists})
}

func (p *Person) PrivateInfo() {
	id := conv.ToInt(p.GetQueryPostForm("id"))
	if id <= 0 {
		p.Json(http.StatusNotImplemented, "illegal parameter")
		return
	}
	var info gin.H
	switch id {
	case 1:
		info = gin.H{"id": 1, "name": "peter", "age": 12, "birthday": "1990-12-11", "address": "Guangzhou, Guangdong Province, China"}
		p.Json(http.StatusOK, "success", gin.H{"info": info})
	default:
		p.Json(http.StatusNoContent, "person does not exist")
	}
}
