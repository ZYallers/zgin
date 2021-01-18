package tool

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ZYallers/zgin/app"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// NowMemStats
func NowMemStats() string {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return fmt.Sprintf("Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes) NumGoroutine:%d", ms.Alloc, ms.HeapIdle, ms.HeapReleased, runtime.NumGoroutine())
}

// PushSimpleMessage
func PushSimpleMessage(msg string, isAtAll bool) {
	host, _ := os.Hostname()
	text := []string{
		msg + "\n---------------------------",
		"App: " + app.Name,
		"Mode: " + gin.Mode(),
		"Listen: " + *app.HttpServerAddr,
		"HostName: " + host,
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"SystemIP: " + SystemIP(),
		"PublicIP: " + PublicIP(),
	}
	if gin.IsDebugging() {
		isAtAll = false // 开发环境下，不需要@所有人，减少干扰!
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": strings.Join(text, "\n") + "\n",
		},
		"at": map[string]interface{}{
			"isAtAll": isAtAll,
		},
	}
	url := "https://oapi.dingtalk.com/robot/send?access_token=" + app.GracefulRobotToken
	_, _ = NewRequest(url).SetHeaders(map[string]string{"Content-Type": "application/json;charset=utf-8"}).SetPostData(postData).Post()
}

// PushContextMessage
func PushContextMessage(ctx *gin.Context, msg string, reqStr string, stack string, isAtAll bool) {
	host, _ := os.Hostname()
	text := []string{
		msg + "\n---------------------------",
		"App: " + app.Name,
		"Mode: " + gin.Mode(),
		"Listen: " + *app.HttpServerAddr,
		"HostName: " + host,
		"Time: " + time.Now().Format("2006/01/02 15:04:05.000"),
		"Url: " + "https://" + ctx.Request.Host + ctx.Request.URL.String(),
		"SystemIP: " + SystemIP(),
		"PublicIP: " + PublicIP(),
		"ClientIP: " + ClientIP(ctx.ClientIP()),
	}
	if reqStr != "" {
		text = append(text, "\nRequest:\n"+strings.ReplaceAll(reqStr, "\n", ""))
	}
	if stack != "" {
		text = append(text, "\nStack:\n"+stack)
	}
	if gin.IsDebugging() {
		isAtAll = false // 开发环境下，不需要@所有人，减少干扰!
	}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": strings.Join(text, "\n") + "\n",
		},
		"at": map[string]interface{}{
			"isAtAll": isAtAll,
		},
	}
	url := "https://oapi.dingtalk.com/robot/send?access_token=" + app.ErrorRobotToken
	_, _ = NewRequest(url).SetHeaders(map[string]string{"Content-Type": "application/json;charset=utf-8"}).SetPostData(postData).Post()
}

// 带签名http请求
// 默认请求方式POST，超时3秒
func HttpRequestWithSign(url string, data map[string]interface{}, params ...interface{}) string {
	var (
		method  = http.MethodPost
		timeout = time.Second * 3
	)
	if len(params) >= 1 {
		method = params[0].(string)
	}
	if len(params) >= 2 {
		timeout = params[1].(time.Duration)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	data["utime"] = timestamp

	hash := md5.New()
	hash.Write([]byte(timestamp + app.TokenKey))
	md5str := hex.EncodeToString(hash.Sum(nil))
	data["sign"] = base64.StdEncoding.EncodeToString([]byte(md5str))

	headers := map[string]string{
		"Connection":   "close",
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "ZGin/1.1.6",
	}
	req := NewRequest(url).SetMethod(method).SetHeaders(headers).SetTimeOut(timeout)

	switch method {
	case http.MethodGet:
		queries := make(map[string]string, len(data))
		for k, v := range data {
			queries[k] = fmt.Sprintf("%v", v)
		}
		req.SetQueries(queries)
	default:
		req.SetPostData(data)
	}

	if resp, err := req.Send(); err == nil {
		return resp.Body
	} else {
		//logger.Use("HttpRequestWithSign").Info(url, zap.String("resp.Body", err.Error()))
		return ""
	}
}
