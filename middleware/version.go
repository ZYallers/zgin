package middleware

import (
	"github.com/ZYallers/zgin/consts"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"github.com/syyongx/php2go"
	"net/http"
)

func VersionCompare(check types.ICheck, handlers types.RestHandlers) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//defer func(t time.Time) { fmt.Println("VersionCompare runtime:", time.Now().Sub(t)) }(time.Now())
		cfgVer, cfgVerKey := check.GetVersion()
		httpVersion := QueryPostForm(ctx, cfgVerKey, cfgVer)
		for _, handler := range handlers {
			if handler.Version.Value == "" || httpVersion == handler.Version.Value {
				SetRestHandler(ctx, handler)
				return
			}
			if handler.Version.Plus {
				if php2go.VersionCompare(httpVersion, handler.Version.Value, ">") {
					SetRestHandler(ctx, handler)
					return
				}
			}
		}
		AbortWithJson(ctx, http.StatusNotFound, "page version mismatch")
	}
}

func SetRestHandler(ctx *gin.Context, handler *types.RestHandler) {
	ctx.Set(consts.RestHandlerKey, handler)
}

func GetRestHandler(ctx *gin.Context) *types.RestHandler {
	if val, ok := ctx.Get(consts.RestHandlerKey); ok && val != nil {
		if handler, ok := val.(*types.RestHandler); ok {
			return handler
		}
	}
	return nil
}
