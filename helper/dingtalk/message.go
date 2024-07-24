package dingtalk

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ZYallers/golib/utils/curl"
	"github.com/ZYallers/zgin/consts"
	"github.com/gin-gonic/gin"
)

var (
	// 连续换行符
	continuousLineBreaksRegExp, _ = regexp.Compile(`\s{2,}`)
)

func PushSimpleMessage(msg interface{}, isAtAll bool) {
	PushMessage(getGracefulToken(), msg, isAtAll)
}

func PushContextMessage(ctx *gin.Context, msg interface{}, reqStr string, stack string, isAtAll bool) {
	PushMessage(getErrorToken(), msg, isAtAll, ctx, reqStr, stack)
}

func PushMessage(token string, msg interface{}, options ...interface{}) {
	defer func() { recover() }()
	title := fmt.Sprint(msg)
	if token == "" || title == "" {
		return
	}
	text := []string{title + "\n---------------------------",
		"App: " + getName(),
		"Mode: " + gin.Mode(),
		"Listen: " + getHttpAddr(),
		"HostName: " + getHostName(),
		"Time: " + time.Now().Format(consts.LogTimeFormat),
		"SystemIP: " + getSystemIP(),
	}
	var (
		optLen  = len(options)
		isAtAll bool
	)
	if optLen > 0 && !gin.IsDebugging() {
		if val, ok := options[0].(bool); ok {
			isAtAll = val
		}
	}
	if optLen > 1 {
		if ctx, ok := options[1].(*gin.Context); ok && ctx != nil {
			text = append(text,
				"ClientIP: "+ctx.ClientIP(),
				"Url: "+ctx.Request.Host+ctx.Request.URL.Path,
			)
		}
	}
	if optLen > 2 {
		if s, ok := options[2].(string); ok && s != "" {
			s := continuousLineBreaksRegExp.ReplaceAllString(s, "\n")
			text = append(text, "\nRequest:\n"+strings.TrimSpace(s))
		}
	}
	if optLen > 3 {
		if s, ok := options[3].(string); ok && s != "" {
			text = append(text, "\nStack:\n"+strings.TrimSpace(s))
		}
	}
	data := map[string]interface{}{
		"msgtype": "text",
		"at":      map[string]interface{}{"isAtAll": isAtAll},
		"text":    map[string]string{"content": strings.Join(text, "\n") + "\n"},
	}
	_, _ = curl.NewRequest("https://oapi.dingtalk.com/robot/send?access_token=" + token).
		SetContentType(consts.JsonContentTypeValue).
		SetTimeOut(3 * time.Second).
		SetPostData(data).
		Post()
}
