package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"strings"
)

var (
	prom    = promhttp.Handler()
	appInfo = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_info",
			Help: "now running go app information.",
		},
		[]string{"name", "cmdline"},
	)
)

func PrometheusHandler(ctx *gin.Context, name string) {
	prometheus.MustRegister(appInfo)
	appInfo.WithLabelValues(name, strings.Join(os.Args, " ")).Inc()
	prom.ServeHTTP(ctx.Writer, ctx.Request)
}
