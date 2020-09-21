package swagger

import (
	"fmt"
	"github.com/ZYallers/zgin/library/tool"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func DocsHandler(ctx *gin.Context) {
	// 生产环境不注册
	if !gin.IsDebugging() {
		return
	}
	pwd, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/doc/swagger.json", pwd)
	if !tool.FileIsExist(filePath) {
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
