package log

import (
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"
)

var log *stdlog.Logger

// fileLog 定义filelog类型，目前仅表示文件名
type fileLog string

func Run(destination string) {
	log = stdlog.New(fileLog(destination), "go: ", stdlog.LstdFlags)
}

// Write fileLog 实现 io.Writer 接口
func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

// RegisterHandlers 注册日志服务
func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
			}
			writeLog(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	})
}

func writeLog(msg string) {
	log.Printf("%v\n", msg)
}
