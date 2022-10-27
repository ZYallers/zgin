package types

import (
	"github.com/ZYallers/golib/funcs/conv"
	"github.com/ZYallers/golib/utils/json"
	"github.com/ZYallers/zgin/consts"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type IController interface {
	SetContext(*gin.Context)
}

type Controller struct {
	ctx *gin.Context
}

func (c *Controller) SetContext(ctx *gin.Context) {
	c.ctx = ctx
}

func (c *Controller) GetContext() *gin.Context {
	return c.ctx
}

func (c *Controller) DumpRequest() string {
	if c.ctx == nil {
		return ""
	}
	s, _ := c.ctx.Get(consts.ReqStrKey)
	return conv.ToString(s)
}

func (c *Controller) GetLoggedUserId() int {
	if c.ctx == nil {
		return 0
	}
	if data, ok := c.ctx.Get(consts.SessDataKey); ok && data != nil {
		if vars, ok := data.(map[string]interface{}); ok {
			if userInfo, ok := vars["userinfo"].(map[string]interface{}); ok {
				if s, ok := userInfo["userid"].(string); ok {
					i, _ := strconv.Atoi(s)
					return i
				}
			}
		}
	}
	return 0
}

func (c *Controller) Json(a ...interface{}) {
	if c.ctx == nil {
		panic("controller context nil pointer")
	}
	var h gin.H
	switch len(a) {
	case 1:
		h = gin.H{"code": conv.ToInt(a[0]), "msg": "", "data": nil}
	case 2:
		h = gin.H{"code": conv.ToInt(a[0]), "msg": conv.ToString(a[1]), "data": nil}
	case 3:
		h = gin.H{"code": conv.ToInt(a[0]), "msg": conv.ToString(a[1]), "data": a[2]}
	case 4:
		h = gin.H{"code": conv.ToInt(a[0]), "msg": conv.ToString(a[1]), "data": a[2], "record": a[3]}
	}

	c.ctx.Abort()
	c.ctx.Header(consts.JsonContentTypeKey, consts.JsonContentTypeValue)
	c.ctx.Status(http.StatusOK)
	if bte, err := json.Marshal(h); err != nil {
		c.ctx.Writer.WriteString(`{"code":500,"msg":"` + conv.ToString(err) + `"}`)
	} else {
		c.ctx.Writer.Write(bte)
	}
}

func (c *Controller) GetQueryPostForm(keys ...string) string {
	if c.ctx == nil {
		return ""
	}
	if len(keys) == 0 {
		return ""
	}
	if val, ok := c.ctx.GetPostForm(keys[0]); ok {
		return val
	}
	if val, ok := c.ctx.GetQuery(keys[0]); ok {
		return val
	}
	if len(keys) == 2 {
		return keys[1]
	}
	return ""
}

func (c *Controller) GetQueryByMethod(key, defaultValue string) string {
	if c.ctx == nil {
		return ""
	}
	var query string
	switch c.ctx.Request.Method {
	case http.MethodPost:
		query = c.ctx.DefaultPostForm(key, defaultValue)
	default:
		query = c.ctx.DefaultQuery(key, defaultValue)
	}
	return query
}

func (c *Controller) QueryPostNumber(key string, defaultValue ...int) int {
	if s := c.GetQueryPostForm(key); s != "" {
		i, _ := strconv.Atoi(s)
		return i
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (c *Controller) Query(keys ...string) string {
	if c.ctx == nil {
		return ""
	}
	if len(keys) == 0 {
		return ""
	}
	if val, ok := c.ctx.GetQuery(keys[0]); ok {
		return val
	}
	if len(keys) == 2 {
		return keys[1]
	}
	return ""
}

func (c *Controller) PostForm(keys ...string) string {
	if c.ctx == nil {
		return ""
	}
	if len(keys) == 0 {
		return ""
	}
	if val, ok := c.ctx.GetPostForm(keys[0]); ok {
		return val
	}
	if len(keys) == 2 {
		return keys[1]
	}
	return ""
}
