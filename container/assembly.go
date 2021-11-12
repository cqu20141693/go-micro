package container

import (
	"encoding/json"
	utils "github.com/cqu20141693/go-micro/utils"
	"github.com/spf13/viper"
)

type ConfigProperties interface {
	prefix() string
}

var ConfigContainer []ConfigProperties = make([]ConfigProperties, 8)
var stop = make(chan bool)

func AddConfigProperties(cpc <-chan ConfigProperties) {
	for {
		select {

		case cp := <-cpc:
			ConfigContainer = append(ConfigContainer, cp)
			marshal, err := json.Marshal(viper.GetString(cp.prefix()))
			if err != nil {
				utils.Logger.Info("config Marshal  failed")
				return
			}
			err = json.Unmarshal(marshal, &cp)
			if err != nil {
				utils.Logger.Info("config Unmarshal failed")
				return
			}
		case <-stop:
			utils.Logger.Info("ConfigProperties add stop")
			break
		}
	}
}

func Reload() {

}
func ConfigUpdate() {
	utils.Logger.Info("config update")
}
func Stop() {
	stop <- true
}
