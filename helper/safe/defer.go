package safe

import (
	"fmt"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/helper/dingtalk"
	"go.uber.org/zap"
	"runtime/debug"
)

func Defer() {
	r := recover()
	if r == nil {
		return
	}
	err := fmt.Sprintf("%v", r)
	stack := string(debug.Stack())
	logger.Use("panic").Error(err, zap.String("debug_stack", stack))
	dingtalk.PushSimpleMessage(fmt.Sprintf("recovery from panic:\n%s\n%s", err, stack), true)
}
