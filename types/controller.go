package types

import (
	"fmt"
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
	reqStr, _ := c.Ctx.Get(consts.ReqStrKey)
	return reqStr.(string)
}

// GetLoggedUserId 获取已登陆用户的ID
func (c *Controller) GetLoggedUserId() int {
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

// Json 输出方法
// args 三个参数: 第一个是code，第二个是msg，第三个是data
func (c *Controller) Json(a ...interface{}) {
	var h gin.H
	switch len(a) {
	case 1:
		h = gin.H{"code": a[0], "msg": "", "data": struct{}{}}
	case 2:
		h = gin.H{"code": a[0], "msg": a[1], "data": struct{}{}}
	case 3:
		h = gin.H{"code": a[0], "msg": a[1], "data": a[2]}
	case 4:
		h = gin.H{"code": a[0], "msg": a[1], "data": a[2], "record": a[3]}
	}
	bte, err := libJson.Marshal(h)
	if err != nil {
		bte = []byte(fmt.Sprintf(`{"code":%d,"msg":"%v"}`, http.StatusInternalServerError, err))
	}
	c.Ctx.Abort()
	c.Ctx.Status(http.StatusOK)
	c.Ctx.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, _ = c.Ctx.Writer.Write(bte)
}

// 从Get和Post里获取Key的值
func (c *Controller) GetQueryPostForm(keys ...string) string {
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

// GetQueryByMethod 根据ctx.httpMethod获取GET或POST参数
func (c *Controller) GetQueryByMethod(key, defaultValue string) string {
	var query string
	switch c.Ctx.Request.Method {
	case http.MethodPost:
		query = c.Ctx.DefaultPostForm(key, defaultValue)
	default:
		query = c.Ctx.DefaultQuery(key, defaultValue)
	}
	return query
}

// QueryPostNumber 从Query或PostForm中获取数字类型key值
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

// Query 从Get里获取Key的值
func (c *Controller) Query(keys ...string) string {
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

// PostForm 从Post里获取Key的值
func (c *Controller) PostForm(keys ...string) string {
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
