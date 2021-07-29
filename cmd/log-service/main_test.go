package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_readConfFile(t *testing.T) {
	conf := readConfFile("./log-service.yaml")
	dataBytes, _ := json.Marshal(conf)
	fmt.Println(string(dataBytes))
}
