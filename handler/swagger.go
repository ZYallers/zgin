package handler

import (
	"fmt"
	"github.com/ZYallers/golib/funcs/files"
	"github.com/ZYallers/zgin/option"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func WithSwagger() option.App {
	return func(app *types.App) {
		app.Server.Http.Handler.(*gin.Engine).GET("/swag/json", SwagHandler)
	}
}

func SwagHandler(ctx *gin.Context) {
	pwd, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/doc/swagger.json", pwd)
	if !files.FileIsExist(filePath) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusNotFound, "msg": "swagger.json file not exist"})
		return
	}
	f, err := os.Open(filePath)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": err.Error()})
		return
	}
	ctx.Header("Content-Type", "application/json;charset=utf-8")
	ctx.String(http.StatusOK, string(fd))
}
