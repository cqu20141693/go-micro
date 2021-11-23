package container

import (
	utils "github.com/cqu20141693/go-micro/utils"
	"syscall"
)

var SingletonFactory map[string]interface{}
var cpchan = make(chan ConfigProperties, 8)
var properties chan ConfigProperties

func init() {
	properties = AddConfigProperties(cpchan)
}
func InjectSingleton(key string, o interface{}) ConfigProperties {
	if _, ok := SingletonFactory[key]; ok {
		utils.Logger.Info("The singleton factory already exists instance=" + key)
		syscall.Exit(1)
	}
	SingletonFactory[key] = o

	if cp, ok := o.(ConfigProperties); ok {
		cpchan <- cp
	}
	return <-properties
}

func Destroy() {
	utils.Logger.Info("container destroy")
}
