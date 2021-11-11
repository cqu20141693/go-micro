package utils

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var ConfigClient config_client.IConfigClient
var NamingClient naming_client.INamingClient
var Err error
var YamlViper *viper.Viper
var Default_cluster = "DEFAULT"
var Default_group = "DEFAULT_GROUP"

var sc = []constant.ServerConfig{
	*constant.NewServerConfig(
		"172.30.203.22",
		8848,
		constant.WithScheme("http"),
		constant.WithContextPath("/nacos"),
	)}

// 创建clientConfig
var cc = constant.ClientConfig{
	NamespaceId:         "ca1d1ded-cb0b-460c-8efa-7e665c7a34e0",
	TimeoutMs:           5000,
	NotLoadCacheAtStart: true,
	LogDir:              "/tmp/nacos/log",
	CacheDir:            "/tmp/nacos/cache",
	RotateTime:          "1h",
	MaxAge:              3,
	LogLevel:            "debug",
}

func init() {
	ConfigClient, Err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if Err != nil {
		panic(Err)
	}
	// 创建服务发现客户端
	NamingClient, Err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if Err != nil {
		panic(Err)
	}
	YamlViper = viper.New()
	YamlViper.SetConfigType("yaml")
}

func GetConfigByDataId(dataId string) (content string, err error) {

	return GetConfig(dataId, Default_group)
}

func GetConfig(dataId, group string) (content string, err error) {

	content, err = ConfigClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	return
}

type ConfigListenHandler func(namespace, group, dataId, data string)

func ListenConfigByDataId(dataId string, handler ConfigListenHandler) error {
	return ListenConfig(dataId, Default_group, handler)
}

func ListenConfig(dataId, group string, handler ConfigListenHandler) error {
	return ConfigClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			handler(namespace, group, dataId, data)
		},
	})
}
