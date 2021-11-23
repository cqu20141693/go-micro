package main

import (
	"fmt"
	"github.com/cqu20141693/go-micro"
	"github.com/cqu20141693/go-micro/container"
	"os"
)

type SipConfig struct {
	Serial        string `json:"serial"`
	Realm         string `json:"realm"`
	ListenAddress string `json:"listenAddress"`
}

func main() {

	micro.Run(os.Args)
	config := SipConfig{}
	container.InjectSingleton("sip", &config)
	fmt.Printf("inject config=%v", config)
}
