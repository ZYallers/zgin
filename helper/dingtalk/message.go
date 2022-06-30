package dingtalk

import (
	"fmt"
	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/golib/utils/curl"
	"github.com/ZYallers/zgin/helper/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"os"
	"strings"
	"time"
)

const (
	timeout   = 3 * time.Second
	uriPrefix = "https://oapi.dingtalk.com/robot/send?access_token="
)

var headers = map[string]string{"Content-Type": "application/json;charset=utf-8"}

func PushSimpleMessage(msg interface{}, isAtAll bool) {
	PushMessage(gracefulToken(), msg, isAtAll)
}

func PushContextMessage(ctx *gin.Context, msg interface{}, reqStr string, stack string, isAtAll bool) {
	PushMessage(errorToken(), msg, isAtAll, ctx, reqStr, stack)
}

func PushMessage(token string, msg interface{}, options ...interface{}) {
	if token == "" || msg == "" {
		return
	}
	defer func() { recover() }()
	text := []string{
		getMsgText(msg) + "\n---------------------------",
		"App: " + appName(),
		"Mode: " + gin.Mode(),
		"Listen: " + appHttpAddr(),
		"HostName: " + hostname(),
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"SystemIP: " + nets.SystemIP(),
		"PublicIP: " + nets.PublicIP(),
	}
	optionsLen := len(options)
	var isAtAll bool
	if !gin.IsDebugging() && optionsLen > 0 {
		if val, ok := options[0].(bool); ok {
			isAtAll = val
		}
	}
	var ctx *gin.Context
	if optionsLen > 1 {
		if val, ok := options[1].(*gin.Context); ok {
			ctx = val
		}
	}
	if ctx != nil {
		text = append(text,
			"ClientIP: "+nets.ClientIP(ctx.ClientIP()),
			"Url: "+"https://"+ctx.Request.Host+ctx.Request.URL.String(),
		)
	}
	if optionsLen > 2 {
		if reqStr, ok := options[2].(string); ok && reqStr != "" {
			text = append(text, "\nRequest:\n"+strings.ReplaceAll(reqStr, "\n", ""))
		}
	}
	if optionsLen > 3 {
		if stack, ok := options[3].(string); ok && stack != "" {
			text = append(text, "\nStack:\n"+stack)
		}
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": strings.Join(text, "\n") + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	_, _ = curl.NewRequest(uriPrefix + token).SetHeaders(headers).SetTimeOut(timeout).SetPostData(postData).Post()
}

func getMsgText(msg interface{}) string {
	var s string
	switch v := msg.(type) {
	case string:
		s = v
	case error:
		s = v.Error()
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

func appName() string {
	if val := config.AppValue("name"); val != nil {
		return cast.ToString(val)
	}
	return "unknown"
}

func appHttpAddr() string {
	if val := config.AppValue("http_addr"); val != nil {
		return cast.ToString(val)
	}
	return "unknown"
}

func hostname() string {
	if val, _ := os.Hostname(); val != "" {
		return val
	}
	return "unknown"
}

func gracefulToken() string {
	if val := config.AppValue("graceful_token"); val != nil {
		return cast.ToString(val)
	}
	return ""
}

func errorToken() string {
	if val := config.AppValue("error_token"); val != nil {
		return cast.ToString(val)
	}
	return ""
}
