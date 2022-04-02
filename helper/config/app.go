package config

import (
	"github.com/spf13/viper"
	"os"
)

const dingtalkDomain = "https://oapi.dingtalk.com/robot/send?access_token="

var (
	appName               string
	appHttpServerAddr     string
	hostname              string
	appGracefulRobotToken string
	appErrorRobotToken    string
)

func Name() string {
	if appName != "" {
		return appName
	}
	if appName = viper.GetString("app.name"); appName == "" {
		return "unknown"
	}
	return appName
}

func HttpAddr() string {
	if appHttpServerAddr != "" {
		return appHttpServerAddr
	}
	if appHttpServerAddr = viper.GetString("app.http_addr"); appHttpServerAddr == "" {
		return "unknown"
	}
	return appHttpServerAddr
}

func Hostname() string {
	if hostname != "" {
		return hostname
	}
	if hostname, _ = os.Hostname(); hostname == "" {
		return "unknown"
	}
	return hostname
}

func GracefulUri() string {
	if appGracefulRobotToken != "" {
		return dingtalkDomain + appGracefulRobotToken
	}
	if appGracefulRobotToken = viper.GetString("app.graceful_token"); appGracefulRobotToken != "" {
		return dingtalkDomain + appGracefulRobotToken
	}
	return ""
}

func ErrorUri() string {
	if appErrorRobotToken != "" {
		return dingtalkDomain + appErrorRobotToken
	}
	if appErrorRobotToken = viper.GetString("app.error_token"); appErrorRobotToken != "" {
		return dingtalkDomain + appErrorRobotToken
	}
	return ""
}
