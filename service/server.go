package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	stdlog "log"
	"net/http"
	"strconv"
	"time"

	"github.com/draymonders/distributed/log"
	"github.com/draymonders/distributed/registry"
)

// Start 启动服务
// 	registerHandlersFunc 服务本身要注册的http方法
func Start(ctx context.Context, host string, port int, service *registry.Service, registerHandlersFunc func()) (context.Context, error) {
	log.SetLog(string(service.Name))
	if service == nil {
		panic("service must init")
	}

	registerHandlersFunc()
	initRegisterToCluster(service)
	ctx = startService(ctx, service.Name, host, port)
	return ctx, nil
}

// initRegisterToCluster 服务第一次注册到注册中心
func initRegisterToCluster(service *registry.Service) {
	registryUrl := fmt.Sprintf("http://%s%s%s", registry.RegisterServiceHost, registry.RegisterServicePort, registry.RegisterServicePath)
	bytesData, err := json.Marshal(service)
	if err != nil {
		panic(fmt.Sprintf("service json.Marshal error: %v", err))
	}
	resp, err := http.Post(registryUrl, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		stdlog.Printf("register to cluster error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		stdlog.Printf("register to cluster fail")
	}
}

// startService 服务启动
func startService(ctx context.Context, serviceName registry.ServiceName, host string, port int) context.Context {
	ctx, cancelFunc := context.WithCancel(ctx)
	url := ":" + strconv.FormatInt(int64(port), 10)

	var srv http.Server
	srv.Addr = url

	go func() {
		// stdlog.Printf("start service: %v url: %v\n", serviceName, url)
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

	return ctx
}
