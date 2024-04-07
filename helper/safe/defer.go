package safe

import (
	"fmt"
	"runtime/debug"

	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/zgin/helper/dingtalk"
	"go.uber.org/zap"
)

func Defer() {
	r := recover()
	if r == nil {
		return
	}

	// 尝试将r转换为error，以获得更好的错误信息
	var err error
	if e, ok := r.(error); ok {
		err = e
	} else {
		err = fmt.Errorf("%v", r)
	}
	stack := string(debug.Stack())
	logger.Use("panic").Error(err.Error(), zap.String("debug_stack", stack))
	dingtalk.PushSimpleMessage(fmt.Sprintf("recovery from panic:\n%s\n%s", err, stack), true)
}
