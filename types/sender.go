package types

import (
	"github.com/ZYallers/zgin/helper/config"
	"github.com/ZYallers/zgin/helper/dingtalk"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

const tokenKey = "sql_token"

type Sender struct {
	token string
}

func (s *Sender) sqlToken() string {
	if s.token == "" {
		if v := config.AppValue(tokenKey); v != nil {
			s.token = cast.ToString(v)
		}
	}
	return s.token
}
func (s *Sender) Open() bool      { return true }
func (s *Sender) Always() bool    { return gin.IsDebugging() }
func (s *Sender) Push(msg string) { dingtalk.PushMessage(s.sqlToken(), msg, true) }
