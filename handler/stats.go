package handler

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
)

func WithStats(acts gin.Accounts) option.App {
	return func(app *types.App) {
		if acts == nil {
			app.Server.Http.Handler.(*gin.Engine).GET("/stats/*filepath", StatsHandler)
			return
		}
		app.Server.Http.Handler.(*gin.Engine).Group("/stats", gin.BasicAuth(acts)).GET("/*filepath", StatsHandler)
	}
}

func StatsHandler(ctx *gin.Context) {
	if ctx.Param("filepath") == "/ws" {
		statsviz.Ws(ctx.Writer, ctx.Request)
		return
	}
	statsviz.IndexAtRoot("/stats").ServeHTTP(ctx.Writer, ctx.Request)
}

// Work loops forever, generating a bunch of allocations of various sizes
// in order to force the garbage collector to work.
func work() {
	m := map[string][]byte{}
	for {
		b := make([]byte, 512+rand.Intn(16*1024))
		m[strconv.Itoa(len(m)%(10*100))] = b

		if len(m)%(10*100) == 0 {
			m = make(map[string][]byte)
		}

		time.Sleep(10 * time.Millisecond)
	}
}
