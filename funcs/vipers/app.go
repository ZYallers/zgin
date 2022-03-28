package vipers

import (
	"github.com/spf13/viper"
	"os"
)

const robotLinkPrefix = "https://oapi.dingtalk.com/robot/send?access_token="

var (
	appName               string
	appHttpServerAddr     string
	hostname              string
	appGracefulRobotToken string
	appErrorRobotToken    string
)

func AppName() string {
	if appName != "" {
		return appName
	}
	if appName = viper.GetString("App.Name"); appName == "" {
		return "unknown"
	}
	return appName
}

func AppHttpServerAddr() string {
	if appHttpServerAddr != "" {
		return appHttpServerAddr
	}
	if appHttpServerAddr = viper.GetString("App.HttpServerAddr"); appHttpServerAddr == "" {
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

func AppGracefulRobotUrl() string {
	if appGracefulRobotToken != "" {
		return robotLinkPrefix + appGracefulRobotToken
	}
	if appGracefulRobotToken = viper.GetString("App.GracefulRobotToken"); appGracefulRobotToken != "" {
		return robotLinkPrefix + appGracefulRobotToken
	}
	return ""
}

func AppErrorRobotUrl() string {
	if appErrorRobotToken != "" {
		return robotLinkPrefix + appErrorRobotToken
	}
	if appErrorRobotToken = viper.GetString("App.ErrorRobotToken"); appErrorRobotToken != "" {
		return robotLinkPrefix + appErrorRobotToken
	}
	return ""
}
