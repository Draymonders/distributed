package util

import "time"

// DoRetry 重试方法
func DoRetry(cnt int, d time.Duration, f func() bool) bool {
	for i := 0; i < cnt; i++ {
		if f() {
			return true
		}
		time.Sleep(d)
	}
	return false
}
