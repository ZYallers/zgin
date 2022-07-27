package dingtalk

import (
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
	PushMessage(getGracefulToken(), msg, isAtAll)
}

func PushContextMessage(ctx *gin.Context, msg interface{}, reqStr string, stack string, isAtAll bool) {
	PushMessage(getErrorToken(), msg, isAtAll, ctx, reqStr, stack)
}

func PushMessage(token string, msg interface{}, options ...interface{}) {
	defer func() { recover() }()
	title := cast.ToString(msg)
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
	_, _ = curl.NewRequest(uriPrefix + token).SetHeaders(headers).SetTimeOut(timeout).SetPostData(data).Post()
}

var (
	systemIP      string
	publicIP      string
	hostName      string
	name          string
	httpAddr      string
	gracefulToken string
	errorToken    string
)

func getHostName() string {
	if hostName == "" {
		if val, _ := os.Hostname(); val != "" {
			hostName = val
		} else {
			return "unknown"
		}
	}
	return hostName
}

func getSystemIP() string {
	if systemIP == "" {
		if s := nets.SystemIP(); s != "" {
			systemIP = s
		} else {
			return "unknown"
		}
	}
	return systemIP
}

func getPublicIP() string {
	if publicIP == "" {
		if s := nets.PublicIP(); s != "" {
			publicIP = s
		} else {
			return "unknown"
		}
	}
	return publicIP
}

func getName() string {
	if name == "" {
		if val := config.AppValue("name"); val != nil {
			name = cast.ToString(val)
		} else {
			return "unknown"
		}
	}
	return name
}

func getHttpAddr() string {
	if httpAddr == "" {
		if val := config.AppValue("http_addr"); val != nil {
			httpAddr = cast.ToString(val)
		} else {
			return "unknown"
		}
	}
	return httpAddr
}

func getGracefulToken() string {
	if gracefulToken == "" {
		if val := config.AppValue("graceful_token"); val != nil {
			gracefulToken = cast.ToString(val)
		}
	}
	return gracefulToken
}

func getErrorToken() string {
	if errorToken == "" {
		if val := config.AppValue("error_token"); val != nil {
			errorToken = cast.ToString(val)
		}
	}
	return errorToken
}
