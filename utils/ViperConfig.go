package micro

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

func ReadLocalConfig() {
	// 读取本地配置
	viper.SetConfigName("bootstrap.yml")
	viper.SetDefault("server.port", 8080)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./resource")
	viper.AddConfigPath("/etc/resource")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Logger.Info("Config file not found; ignore error if desired")

		} else {
			Logger.Info("Config file was found but another error was produced")
		}
		log.Fatal(err)
	}
}

func getAppConfigName(name, active string) (string, error) {
	if name == "" {
		return "", errors.New("cc.application.name not config")
	} else if active == "" {
		return "", errors.New("cc.profiles.active not config")
	}
	return name + "-" + active + "." + LocalNacosConfig.FileExtension, nil
}
