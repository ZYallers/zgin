package dingtalk

import (
	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/zgin/helper/config"
	"github.com/spf13/cast"
	"os"
)

var systemIP, publicIP, hostName, name, httpAddr, gracefulToken, errorToken string

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
