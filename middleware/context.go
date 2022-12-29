package middleware

import (
	"github.com/ZYallers/golib/funcs/conv"
	"github.com/ZYallers/golib/utils/json"
	"github.com/ZYallers/zgin/consts"
	"github.com/gin-gonic/gin"
	"net/http"
)

func QueryPostForm(ctx *gin.Context, keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	if val, ok := ctx.GetPostForm(keys[0]); ok {
		return val
	}
	if val, ok := ctx.GetQuery(keys[0]); ok {
		return val
	}
	if len(keys) == 2 {
		return keys[1]
	}
	return ""
}

func AbortWithJson(ctx *gin.Context, a ...interface{}) {
	var h gin.H
	switch len(a) {
	case 1:
		h = gin.H{"code": a[0], "msg": ""}
	case 2:
		h = gin.H{"code": a[0], "msg": a[1]}
	case 3:
		h = gin.H{"code": a[0], "msg": a[1], "data": a[2]}
	}
	ctx.Abort()
	ctx.Header(consts.JsonContentTypeKey, consts.JsonContentTypeValue)
	ctx.Status(http.StatusOK)
	if bte, err := json.Marshal(h); err != nil {
		_, _ = ctx.Writer.WriteString(`{"code":500,"msg":"` + conv.ToString(err) + `"}`)
	} else {
		_, _ = ctx.Writer.Write(bte)
	}
}
