// Copyright (c) 2020 HXS R&D Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
//
// @Title api
// @Description Description
//
// @Author zhongyongbiao 2020/9/21 上午11:19
// @Version 1.0.0
// @Software GoLand
package restful

import (
	v100Test "github.com/ZYallers/zgin/controller/v100/test"
	"github.com/ZYallers/zgin/library/expvar"
	"github.com/ZYallers/zgin/library/prometheus"
	"github.com/ZYallers/zgin/library/swagger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Api = &Rest{
	// v1.0.0
	"expvar":    {{Version: "1.0.0+", Method: map[string]byte{http.MethodGet: 1}, Handler: expvar.RunningStatsHandler}},
	"metrics":   {{Version: "1.0.0+", Method: map[string]byte{http.MethodGet: 1}, Handler: prometheus.ServerHandler}},
	"swag/json": {{Version: "1.0.0+", Method: map[string]byte{http.MethodGet: 1}, Handler: swagger.DocsHandler}},
	"test/isok": {{Version: "1.0.0+", Method: map[string]byte{http.MethodGet: 1}, Handler: func(c *gin.Context) { v100Test.Index(c).CheckOk() }}},
}