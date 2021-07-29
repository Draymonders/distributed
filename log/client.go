package log

import (
	"fmt"
	stdlog "log"
)

func SetLog(serviceName string) {
	stdlog.SetPrefix(fmt.Sprintf("[%s] ", serviceName))
	stdlog.SetFlags(stdlog.LstdFlags)
}
