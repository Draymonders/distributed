package util

import (
	stdlog "log"
	"runtime"
)

func recoverFunc() func() {
	return func() {
		// err 就是 panic的内容
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			stdlog.Printf("panic info %v\n", buf)
		}
	}
}
