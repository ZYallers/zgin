package dingtalk

import (
	"os"

	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/zgin/helper/config"
)

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
		if v := config.AppValue("name"); v != nil {
			name = v.(string)
		} else {
			return "unknown"
		}
	}
	return name
}

func getHttpAddr() string {
	if httpAddr == "" {
		if v := config.AppValue("http_addr"); v != nil {
			httpAddr = v.(string)
		} else {
			return "unknown"
		}
	}
	return httpAddr
}

func getGracefulToken() string {
	if gracefulToken == "" {
		if v := config.AppValue("graceful_token"); v != nil {
			gracefulToken = v.(string)
		}
	}
	return gracefulToken
}

func getErrorToken() string {
	if errorToken == "" {
		if v := config.AppValue("error_token"); v != nil {
			errorToken = v.(string)
		}
	}
	return errorToken
}
