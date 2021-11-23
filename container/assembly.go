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
var output = make(chan ConfigProperties)

func AddConfigProperties(cpc <-chan ConfigProperties) chan ConfigProperties {
	go func() {
		for {
			select {

			case cp := <-cpc:
				handleBind(cp)
			case <-stop:
				utils.Logger.Info("ConfigProperties add stop")
				break
			}
		}
	}()
	return output
}

func handleBind(cp ConfigProperties) {
	ConfigContainer = append(ConfigContainer, cp)
	marshal, err := json.Marshal(viper.GetString(cp.prefix()))
	if err != nil {
		utils.Logger.Info("config Marshal  failed")

	}
	err = json.Unmarshal(marshal, &cp)
	if err != nil {
		utils.Logger.Info("config Unmarshal failed")
	}
	output <- cp
}

func Reload() {

}
func ConfigUpdate() {
	utils.Logger.Info("config update")
}
func Stop() {
	stop <- true
}
