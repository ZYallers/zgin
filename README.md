# ZGin API Framework
[![Go Report Card](https://goreportcard.com/badge/github.com/ZYallers/zgin)](https://goreportcard.com/report/github.com/ZYallers/zgin)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/ZYallers/zgin.svg?branch=master)](https://travis-ci.org/ZYallers/zgin) 
[![Foundation](https://img.shields.io/badge/Golang-Foundation-green.svg)](http://golangfoundation.org) 
[![GoDoc](https://pkg.go.dev/badge/github.com/ZYallers/zgin?status.svg)](https://pkg.go.dev/github.com/ZYallers/zgin?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/ZYallers/zgin/-/badge.svg)](https://sourcegraph.com/github.com/ZYallers/zgin?badge)
[![Release](https://img.shields.io/github/release/ZYallers/zgin.svg?style=flat-square)](https://github.com/ZYallers/zgin/releases)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/ZYallers/zgin)](https://www.tickgit.com/browse?repo=github.com/ZYallers/zgin)
[![goproxy.cn](https://goproxy.cn/stats/github.com/ZYallers/zgin/badges/download-count.svg)](https://goproxy.cn)

ZGin is a API framework written in Go (Golang). 
An MVCS, Restful, and version control framework based on the [Gin](https://github.com/gin-gonic/gin) framework.
If you need performance and good productivity, you will love ZGin.


## Installation
To install zgin package, you need to install Go and set your Go workspace first.

1. The first need Go installed (version 1.11+ is required), then you can use the below Go command to install zgin.
```bash
$ go get -u github.com/ZYallers/zgin
```

2. Import it in your code:
```go 
import "github.com/ZYallers/zgin" 
```

## Quick start
```go
# assume the following codes in example.go file
$ cat zgin_test.go
```
```go
package zgin

import (
	"github.com/ZYallers/zgin/app"
	v000T "github.com/ZYallers/zgin/controller/v000/test"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
	"time"
)

var Api = &app.Restful{
	"test/isok": {{Method: map[string]byte{http.MethodGet: 1}, Handler: func(c *gin.Context) { v000T.Index(c).CheckOk() }}},
}

func TestServer(t *testing.T) {
	EnvInit()
	MiddlewareGlobalRegister()
	ExpVarRegister()
	PrometheusRegister()
	if gin.IsDebugging() {
		SwaggerRegister()
		PProfRegister()
		app.HttpServerWriteTimeout = time.Minute
	}
	SessionClientRegister(nil)
	MiddlewareCustomRegister(Api)
	ListenAndServe(app.HttpServerReadTimeout, app.HttpServerWriteTimeout, app.HttpServerShutDownTimeout)
}
```
test main_test.go and visit 0.0.0.0:9010/test/isok (for windows "localhost:8080/test/isok") on browser
```bash
$ go test -v -run TestServer main_test.go
```
## ZGin v1. stable
- MVCS four-tier architecture support
- Restful interface style support
- API version control and permission custom configuration
- PProf middleware support
- Prometheus middleware support
- Swagger api docs middleware support
- Graceful server shutdown and reload

## Build with jsoniter
ZGin uses encoding/json as default json package but you can change to [jsoniter](https://github.com/json-iterator/go) by build from other tags.
```bash
$ go build -tags=jsoniter .
```

## API Examples
```go
var Api = &app.Restful{
	"test/isok": {{Method: map[string]byte{http.MethodGet: 1}, Handler: func(c *gin.Context) { v000T.Index(c).CheckOk() }}},
}
```

## License
Released under the [MIT License](https://github.com/ZYallers/zgin/blob/master/LICENSE)



