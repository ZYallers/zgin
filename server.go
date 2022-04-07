package zgin

import (
	"github.com/ZYallers/golib/funcs/nets"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/helper/dingtalk"
	"github.com/ZYallers/zgin/helper/grace"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutDownTimeout time.Duration
	Http            *http.Server
}

func (s *Server) HealthHandler() *Server {
	s.Http.Handler.(*gin.Engine).GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, `ok`)
	})
	return s
}

func (s *Server) NoRouteHandler() *Server {
	s.Http.Handler.(*gin.Engine).NoRoute(func(ctx *gin.Context) {
		reqStr := ctx.GetString(consts.ReqStrKey)
		p := ctx.Request.URL.Path
		logger.Use("404").Info(p,
			zap.String("proto", ctx.Request.Proto),
			zap.String("method", ctx.Request.Method),
			zap.String("host", ctx.Request.Host),
			zap.String("url", ctx.Request.URL.String()),
			zap.String("query", ctx.Request.URL.RawQuery),
			zap.String("clientIP", nets.ClientIP(ctx.ClientIP())),
			zap.Any("header", ctx.Request.Header),
			zap.String("request", reqStr),
		)
		dingtalk.PushContextMessage(ctx, strings.TrimLeft(p, "/")+" page not found", reqStr, "", false)
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "msg": "page not found"})
	})
	return s
}

func (s *Server) Start() {
	grace.Graceful(s.Http, s.ShutDownTimeout)
}
