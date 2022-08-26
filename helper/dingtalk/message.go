package dingtalk

import (
	"fmt"
	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/golib/utils/curl"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

const (
	timeout = 3 * time.Second
	prefix  = "https://oapi.dingtalk.com/robot/send?access_token="
)

var headers = map[string]string{"Content-Type": "application/json;charset=utf-8"}

func PushSimpleMessage(msg interface{}, isAtAll bool) {
	PushMessage(getGracefulToken(), msg, isAtAll)
}

func PushContextMessage(ctx *gin.Context, msg interface{}, reqStr string, stack string, isAtAll bool) {
	PushMessage(getErrorToken(), msg, isAtAll, ctx, reqStr, stack)
}

func PushMessage(token string, msg interface{}, options ...interface{}) {
	defer func() { recover() }()
	title := fmt.Sprintf("%v", msg)
	if token == "" || title == "" {
		return
	}
	text := []string{title + "\n---------------------------",
		"App: " + getName(),
		"Mode: " + gin.Mode(),
		"Listen: " + getHttpAddr(),
		"HostName: " + getHostName(),
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"SystemIP: " + getSystemIP(),
		"PublicIP: " + getPublicIP(),
	}
	ol := len(options)
	var isAtAll bool
	if ol > 0 && !gin.IsDebugging() {
		if val, ok := options[0].(bool); ok {
			isAtAll = val
		}
	}
	if ol > 1 {
		if ctx, ok := options[1].(*gin.Context); ok && ctx != nil {
			text = append(text, "ClientIP: "+nets.ClientIP(ctx.ClientIP()), "Url: "+"https://"+ctx.Request.Host+ctx.Request.URL.String())
		}
	}
	if ol > 2 {
		if rs, ok := options[2].(string); ok && rs != "" {
			text = append(text, "\nRequest:\n"+strings.ReplaceAll(rs, "\n", ""))
		}
	}
	if ol > 3 {
		if stack, ok := options[3].(string); ok && stack != "" {
			text = append(text, "\nStack:\n"+stack)
		}
	}
	data := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": strings.Join(text, "\n") + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	_, _ = curl.NewRequest(prefix + token).SetHeaders(headers).SetTimeOut(timeout).SetPostData(data).Post()
}
