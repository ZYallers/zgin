package mvcs

import (
	"github.com/ZYallers/zgin/consts"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type IController interface {
	SetContext(*gin.Context)
}

// Controller
type Controller struct {
	Ctx *gin.Context
}

// SetContext
func (c *Controller) SetContext(ctx *gin.Context) {
	c.Ctx = ctx
}

// DumpRequest
func (c *Controller) DumpRequest() string {
	reqStr, _ := c.Ctx.Get(consts.ReqStrKey)
	return reqStr.(string)
}

// Json 输出方法
// args 三个参数：
// 第一个是code int，代表状态码
// 第二个是msg string，代表信息
// 第三个是data gin.H，代表数据
func (c *Controller) Json(args ...interface{}) {
	switch len(args) {
	case 1:
		c.Ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": args[0], "msg": "ok", "data": nil})
	case 2:
		c.Ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": args[0], "msg": args[1], "data": nil})
	case 3:
		c.Ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": args[0], "msg": args[1], "data": args[2]})
	default:
		c.Ctx.AbortWithStatusJSON(http.StatusOK, args)
	}
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

// 获取已登陆用户的ID
func (c *Controller) GetLoggedUserId() int {
	if data, ok := c.Ctx.Get(consts.SessDataKey); ok && data != nil {
		vars := data.(map[string]interface{})
		if userInfo, ok := vars["userinfo"].(map[string]interface{}); ok {
			if str, ok := userInfo["userid"].(string); ok && str != "" {
				userId, _ := strconv.Atoi(str)
				return userId
			}
		}
	}
	return 0
}

// 根据ctx.httpMethod获取GET或POST参数
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
