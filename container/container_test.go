package container

import (
	"github.com/cqu20141693/go-micro"
	"testing"
)

type SipConfig struct {
	Serial        string `json:"serial"`
	Realm         string `json:"realm"`
	ListenAddress string `json:"listenAddress"`
}

func TestConfigUpdate(t *testing.T) {

	micro.Run([]string{"test"})
	config := SipConfig{}
	InjectSingleton("sip", &config)

}
