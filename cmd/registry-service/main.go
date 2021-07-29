package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"time"

	"github.com/draymonders/distributed/log"
	"github.com/draymonders/distributed/registry"
)

func main() {
	const svcName = "register-service"

	ctx, cancelFunc := context.WithCancel(context.Background())

	registry.InitRegistry()
	registry.RegisterHandlers()
	// 启动http 服务
	srv := http.Server{}
	srv.Addr = registry.RegisterServicePort
	log.SetLog(svcName)

	go func() {
		stdlog.Println(srv.ListenAndServe())
		_ = srv.Shutdown(ctx)
		cancelFunc()
	}()
	time.Sleep(100 * time.Millisecond)

	go func() {
		stdlog.Printf("service started, press to kill\n")
		var content string
		n, _ := fmt.Scanln(&content)
		// 读取命令行内容
		for n == 0 {
			n, _ = fmt.Scanln(&content)
		}
		cancelFunc()
	}()
	<-ctx.Done()
}
