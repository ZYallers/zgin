package v110

import (
	"github.com/ZYallers/zgin/libraries/mvcs"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Index struct {
	mvcs.Controller
	tag struct {
		Two   func() `path:"test/two" ver:"1.1.0+" http:"get,post" sort:"2"`
		Third func() `path:"test/third" ver:"1.1.0+" http:"get,post" sign:"on"`
	}
}

func (i *Index) Two() {
	i.Json(http.StatusOK, "ok", gin.H{"name": "v110.Index.Two"})
}

func (i *Index) Third() {
	i.Json(http.StatusOK, "ok", gin.H{"name": "v110.Index.Third"})
}
