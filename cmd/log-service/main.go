package main

import (
	"context"
	"io/ioutil"
	stdlog "log"
	"os"

	"github.com/draymonders/distributed/registry"
	"gopkg.in/yaml.v2"

	"github.com/draymonders/distributed/log"
	"github.com/draymonders/distributed/service"
)

func main() {
	ctx := context.Background()
	const confFile = "./log-service.yaml"
	conf := readConfFile(confFile)

	serviceName := conf.Name
	host := conf.Host
	port := conf.Port
	svc := registry.NewService(serviceName, host, port)

	log.Run(conf.LogFile)
	ctx, err := service.Start(ctx, host, port, svc, log.RegisterHandlers)
	if err != nil {
		stdlog.Fatalf(err.Error())
		os.Exit(1)
	}
	<-ctx.Done()
}

func readConfFile(filePath string) *Conf {
	dataBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		stdlog.Printf("read conf file error: %v", err)
		panic("read conf file error")
	}
	conf := &Conf{}
	err = yaml.Unmarshal(dataBytes, conf)
	if err != nil {
		stdlog.Printf("yaml.Unmarshal error: %v", err)
		panic("parse data struct error")
	}
	return conf
}

type Conf struct {
	LogServiceConf `yaml:"LogService"`
}

type LogServiceConf struct {
	Name    string `yaml:"Name"`
	Host    string `yaml:"Host"`
	Port    int    `yaml:"Port"`
	LogFile string `yaml:"LogFile"`
}
