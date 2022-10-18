package types

import (
	"fmt"
	"github.com/ZYallers/golib/funcs/conv"
	libJson "github.com/ZYallers/golib/utils/json"
	"github.com/ZYallers/zgin/consts"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type IController interface {
	SetContext(*gin.Context)
}

type Controller struct {
	Ctx *gin.Context
}

func (c *Controller) SetContext(ctx *gin.Context) {
	c.Ctx = ctx
}

func (c *Controller) DumpRequest() string {
	if c.Ctx == nil {
		return ""
	}
	s, _ := c.Ctx.Get(consts.ReqStrKey)
	return conv.ToString(s)
}

func (c *Controller) GetLoggedUserId() int {
	if c.Ctx == nil {
		return 0
	}
	if data, ok := c.Ctx.Get(consts.SessDataKey); ok && data != nil {
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
	if c.Ctx == nil {
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
	c.Ctx.Status(http.StatusOK)
	c.Ctx.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")
	if bte, err := libJson.Marshal(h); err != nil {
		s := fmt.Sprintf(`{"code":%d,"msg":"%v","data":null}`, http.StatusInternalServerError, err)
		_, _ = c.Ctx.Writer.WriteString(s)
	} else {
		_, _ = c.Ctx.Writer.Write(bte)
	}
	c.Ctx.Abort()
}

func (c *Controller) GetQueryPostForm(keys ...string) string {
	if c.Ctx == nil {
		return ""
	}
	if len(keys) == 0 {
		return ""
	}
	if val, ok := c.Ctx.GetQuery(keys[0]); ok {
		return val
	}
	if val, ok := c.Ctx.GetPostForm(keys[0]); ok {
		return val
	}
	if len(keys) == 2 {
		return keys[1]
	}
	return ""
}

func (c *Controller) GetQueryByMethod(key, defaultValue string) string {
	if c.Ctx == nil {
		return ""
	}
	var query string
	switch c.Ctx.Request.Method {
	case http.MethodPost:
		query = c.Ctx.DefaultPostForm(key, defaultValue)
	default:
		query = c.Ctx.DefaultQuery(key, defaultValue)
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
	if c.Ctx == nil {
		return ""
	}
	if len(keys) == 0 {
		return ""
	}
	if val, ok := c.Ctx.GetQuery(keys[0]); ok {
		return val
	}
	if len(keys) == 2 {
		return keys[1]
	}
	return ""
}

func (c *Controller) PostForm(keys ...string) string {
	if c.Ctx == nil {
		return ""
	}
	if len(keys) == 0 {
		return ""
	}
	if val, ok := c.Ctx.GetPostForm(keys[0]); ok {
		return val
	}
	if len(keys) == 2 {
		return keys[1]
	}
	return ""
}
