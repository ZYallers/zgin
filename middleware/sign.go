package middleware

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/ZYallers/zgin/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func SignCheck(check types.ICheck) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//defer func(t time.Time) { fmt.Println("SignCheck runtime:", time.Now().Sub(t)) }(time.Now())
		if handler := GetRestHandler(ctx); handler == nil {
			AbortWithJson(ctx, http.StatusForbidden, "signature handler not found")
			return
		} else {
			if handler.Sign {
				if !verifySign(ctx, check) {
					AbortWithJson(ctx, http.StatusForbidden, "signature verification failed")
				}
			}
		}
	}
}

func verifySign(ctx *gin.Context, check types.ICheck) bool {
	secretKey, key, timeKey, dev, expiration := check.GetSign()
	sign := QueryPostForm(ctx, key)
	if sign == "" {
		return false
	}
	sign, _ = url.QueryUnescape(sign)
	if gin.IsDebugging() && sign == dev {
		return true
	}
	times := QueryPostForm(ctx, timeKey)
	if times == "" {
		return false
	}
	timestamp, err := strconv.ParseInt(times, 10, 0)
	if err != nil {
		return false
	}
	if time.Now().Unix()-timestamp > expiration {
		return false
	}
	hash := md5.New()
	hash.Write([]byte(times + secretKey))
	md5str := hex.EncodeToString(hash.Sum(nil))
	if sign == base64.StdEncoding.EncodeToString([]byte(md5str)) {
		return true
	}
	return false
}
