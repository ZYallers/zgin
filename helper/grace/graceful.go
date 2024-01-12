package grace

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZYallers/zgin/helper/dingtalk"
	"go.uber.org/zap"
)

func pushLog(msg string, fn func(string, ...zap.Field)) {
	fn(msg)
	dingtalk.PushSimpleMessage(msg, true)
}

func Graceful(srv *http.Server, timeout time.Duration, logger *zap.Logger) {
	go func() {
		pushLog(fmt.Sprintf("server(%d) is ready to listen and serve", syscall.Getpid()), logger.Info)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pushLog(fmt.Sprintf("server listen and serve error: %v", err), logger.Error)
			os.Exit(1)
		}
	}()

	quitChan := make(chan os.Signal, 1)
	// SIGTERM 结束程序(kill pid)(可以被捕获、阻塞或忽略)
	// SIGHUP 终端控制进程结束(终端连接断开)
	// SIGINT 用户发送INTR字符(Ctrl+C)触发
	// SIGQUIT 用户发送QUIT字符(Ctrl+/)触发
	signal.Notify(quitChan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	sign := <-quitChan

	// 保证quitChan将不再接收信号
	signal.Stop(quitChan)

	// 控制是否启用HTTP保持活动，默认情况下始终启用保持活动，只有资源受限的环境或服务器在关闭过程中才应禁用它们
	srv.SetKeepAlivesEnabled(false)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pid := syscall.Getpid()
	pushLog(fmt.Sprintf("server(%d) is shutting down(%v)...", pid, sign), logger.Info)
	if err := srv.Shutdown(ctx); err != nil {
		pushLog(fmt.Sprintf("server gracefully shutdown error: %v", err), logger.Error)
	} else {
		pushLog(fmt.Sprintf("server(%d) has stopped", pid), logger.Info)
	}
}
