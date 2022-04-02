package runtimes

import (
	"fmt"
	"runtime"
)

func NowMemStats() string {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return fmt.Sprintf("Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes) NumGoroutine:%d",
		ms.Alloc, ms.HeapIdle, ms.HeapReleased, runtime.NumGoroutine())
}
