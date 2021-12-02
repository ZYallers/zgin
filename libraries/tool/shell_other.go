// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!plan9,!solaris

package tool

import (
	"context"
	"errors"
)

// 执行shell命令，可设置执行超时时间
func ExecShellWithContext(ctx context.Context, command string) (string, error) {
	return "", errors.New("this function does not support running under the current system")
}
