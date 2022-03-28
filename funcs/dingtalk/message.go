package dingtalk

import (
	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/golib/utils/curl"
	"github.com/ZYallers/zgin/funcs/vipers"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

var (
	headers = map[string]string{"Content-Type": "application/json;charset=utf-8"}
)

func PushSimpleMessage(msg string, isAtAll bool) {
	text := []string{
		msg + "\n---------------------------",
		"App: " + vipers.AppName(),
		"Mode: " + gin.Mode(),
		"Listen: " + vipers.AppHttpServerAddr(),
		"HostName: " + vipers.Hostname(),
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"SystemIP: " + nets.SystemIP(),
		"PublicIP: " + nets.PublicIP(),
	}
	if gin.IsDebugging() {
		isAtAll = false
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": strings.Join(text, "\n") + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	_, _ = curl.NewRequest(vipers.AppGracefulRobotUrl()).SetHeaders(headers).
		SetTimeOut(3 * time.Second).SetPostData(postData).Post()
}

func PushContextMessage(ctx *gin.Context, msg string, reqStr string, stack string, isAtAll bool) {
	text := []string{
		msg + "\n---------------------------",
		"App: " + vipers.AppName(),
		"Mode: " + gin.Mode(),
		"Listen: " + vipers.AppHttpServerAddr(),
		"HostName: " + vipers.Hostname(),
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"Url: " + "https://" + ctx.Request.Host + ctx.Request.URL.String(),
		"SystemIP: " + nets.SystemIP(),
		"PublicIP: " + nets.PublicIP(),
		"ClientIP: " + nets.ClientIP(ctx.ClientIP()),
	}
	if reqStr != "" {
		text = append(text, "\nRequest:\n"+strings.ReplaceAll(reqStr, "\n", ""))
	}
	if stack != "" {
		text = append(text, "\nStack:\n"+stack)
	}
	if gin.IsDebugging() {
		isAtAll = false
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": strings.Join(text, "\n") + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	_, _ = curl.NewRequest(vipers.AppErrorRobotUrl()).SetHeaders(headers).
		SetTimeOut(3 * time.Second).SetPostData(postData).Post()
}
