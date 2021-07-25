package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Start 启动服务
func Start(ctx context.Context, serviceName, host, port string, registerHandlersFunc func()) (context.Context, error) {
	// 注册对应的 http handler方法
	registerHandlersFunc()

	ctx = startService(ctx, serviceName, host, port)
	return ctx, nil
}

func startService(ctx context.Context, serviceName, host, port string) context.Context {
	ctx, cancelFunc := context.WithCancel(ctx)
	url := ":" + port

	var srv http.Server
	srv.Addr = url

	go func() {
		log.Printf("start service: %v url: %v\n", serviceName, url)
		log.Println(srv.ListenAndServe())
		cancelFunc()
	}()

	time.Sleep(100 * time.Millisecond)
	go func() {
		log.Printf("press the content to kill the service: %v\n", serviceName)
		var content string
		n, _ := fmt.Scanln(&content)
		// 读取命令行内容
		for n == 0 {
			n, _ = fmt.Scanln(&content)
		}
		cancelFunc()
	}()

	return ctx
}
