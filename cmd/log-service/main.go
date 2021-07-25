package main

import (
	"context"
	stdlog "log"
	"os"

	"github.com/draymonders/distributed/log"
	"github.com/draymonders/distributed/service"
)

func main() {
	ctx := context.Background()
	serviceName := "log-service"
	host, port := "localhost", "4000"

	log.Run("./distributed.log")
	ctx, err := service.Start(ctx, serviceName, host, port, log.RegisterHandlers)
	if err != nil {
		stdlog.Fatalf(err.Error())
		os.Exit(1)
	}
	<-ctx.Done()
}
