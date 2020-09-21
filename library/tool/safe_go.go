package tool

import (
	"fmt"
	"github.com/ZYallers/go-frame/gin/library/logger"
	"go.uber.org/zap"
	"reflect"
	"runtime/debug"
	"sync/atomic"
)

func SafeDefer(params ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%+v", r)
			stack := string(debug.Stack())
			logger.Use("panic").Error(msg, zap.String("debug_stack", stack))
			PushSimpleMessage(fmt.Sprintf("recovery from panic:\n%s\n%s", msg, stack), true)
			return
		}
	}()

	r := recover()
	if r == nil {
		return
	}

	err := fmt.Errorf("%+v", r)
	stack := string(debug.Stack())
	logger.Use("panic").Error(err.Error(), zap.String("debug_stack", stack))
	PushSimpleMessage(fmt.Sprintf("recovery from panic:\n%s\n%s", err.Error(), stack), true)

	if paramsLen := len(params); paramsLen > 0 {
		if reflect.TypeOf(params[0]).String()[0:4] != "func" {
			return
		}
		var args []reflect.Value
		if paramsLen > 1 {
			args = append(args, reflect.ValueOf(err))
			for _, v := range params[1:] {
				args = append(args, reflect.ValueOf(v))
			}
		}
		reflect.ValueOf(params[0]).Call(args)
	}
}

func SafeGo(params ...interface{}) {
	if len(params) == 0 {
		return
	}

	pg := &panicGroup{panics: make(chan string, 1), dones: make(chan int, 1)}
	defer pg.closeChan()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				pg.panics <- fmt.Sprintf("%+v\n%s", r, string(debug.Stack()))
				return
			}
			pg.dones <- 1
		}()
		var args []reflect.Value
		if len(params) > 1 {
			for _, v := range params[1:] {
				args = append(args, reflect.ValueOf(v))
			}
		}
		reflect.ValueOf(params[0]).Call(args)
	}()

	for {
		select {
		case <-pg.dones:
			return
		case p := <-pg.panics:
			panic(p)
		}
	}
}

func NewSafeGo() *panicGroup {
	return &panicGroup{
		panics: make(chan string, 8),
		dones:  make(chan int, 8),
	}
}

type panicGroup struct {
	panics chan string // 协程 panic 通知通道
	dones  chan int    // 协程完成通知通道
	jobN   int32       // 协程并发数量
}

func (g *panicGroup) Go(params ...interface{}) *panicGroup {
	if len(params) == 0 {
		params = []interface{}{func() {}}
	}
	atomic.AddInt32(&g.jobN, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.panics <- fmt.Sprintf("%+v\n%s", r, string(debug.Stack()))
				return
			}
			g.dones <- 1
		}()
		var args []reflect.Value
		if len(params) > 1 {
			for _, v := range params[1:] {
				args = append(args, reflect.ValueOf(v))
			}
		}
		reflect.ValueOf(params[0]).Call(args)
	}()
	return g
}

func (g *panicGroup) Wait() {
	defer g.closeChan()
	for {
		select {
		case <-g.dones:
			if atomic.AddInt32(&g.jobN, -1) == 0 {
				return
			}
		case p := <-g.panics:
			panic(p)
		}
	}
}

func (g *panicGroup) closeChan() {
	close(g.dones)
	close(g.panics)
}
