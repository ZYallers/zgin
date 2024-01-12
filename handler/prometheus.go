package handler

import (
	"os"
	"strings"
	"sync"

	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	once    sync.Once
	prom    = promhttp.Handler()
	appInfo = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_info",
			Help: "now running go app information.",
		},
		[]string{"name", "cmdline"},
	)
)

func WithPrometheus() option.App {
	return func(app *types.App) {
		app.Server.Http.Handler.(*gin.Engine).GET("/metrics", func(ctx *gin.Context) {
			PrometheusHandler(ctx, app.Name)
		})
	}
}

func PrometheusHandler(ctx *gin.Context, name string) {
	once.Do(func() {
		prometheus.MustRegister(appInfo)
		appInfo.WithLabelValues(name, strings.Join(os.Args, " ")).Inc()
	})
	prom.ServeHTTP(ctx.Writer, ctx.Request)
}
