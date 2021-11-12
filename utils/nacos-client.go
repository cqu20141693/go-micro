package micro

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
)

var ConfigClient config_client.IConfigClient
var NamingClient naming_client.INamingClient
var Err error
var DefaultCluster = "DEFAULT"
var DefaultGroup = "DEFAULT_GROUP"

var sc []constant.ServerConfig

func createServerConfig(ipAddr string, port uint64) constant.ServerConfig {
	return *constant.NewServerConfig(
		ipAddr,
		port,
		constant.WithScheme("http"),
		constant.WithContextPath("/nacos"),
	)
}

// 创建clientConfig
var cc constant.ClientConfig

type config struct {
	ServerAddr      string `json:"server-addr"`
	FileExtension   string `json:"file-extension"`
	Namespace       string `json:"namespace"`
	ClusterName     string `json:"cluster-name"`
	Group           string `json:"group"`
	RegisterEnabled bool   `json:"register-enabled"`
}

func newConfig(serverAddr string, fileExtension string, namespace string, clusterName string, group string, registerEnabled bool) *config {
	return &config{ServerAddr: serverAddr, FileExtension: fileExtension, Namespace: namespace, ClusterName: clusterName, Group: group, RegisterEnabled: registerEnabled}
}
func newDefaultConfig() config {
	return config{ServerAddr: "localhost:8848", FileExtension: "yml", Namespace: "", ClusterName: "DEFAULT", Group: "DEFAULT_GROUP", RegisterEnabled: true}
}

var LocalNacosConfig = newDefaultConfig()

func NacosInit() {
	ReadLocalConfig()

	nacosConfig := viper.GetStringMap("cc.cloud.nacos.config")
	Logger.Debug(ToJSONString(nacosConfig))
	marshal, _ := json.Marshal(nacosConfig)
	json.Unmarshal(marshal, &LocalNacosConfig)
	splits := strings.Split(LocalNacosConfig.ServerAddr, ",")
	for i := range splits {
		host := strings.Split(splits[i], ":")
		if len(host) == 2 {
			port, err := strconv.Atoi(host[1])
			if err != nil {
				log.Fatal(err)
			}
			sc = append(sc, createServerConfig(host[0], uint64(port)))
		}
	}
	cc = constant.ClientConfig{
		NamespaceId: LocalNacosConfig.Namespace,
	}

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
	ReadNacosConfig()
}

func ReadNacosConfig() {
	// 读取nacos配置
	name := viper.GetString("cc.application.name")
	active := viper.GetString("cc.profiles.active")
	appConfigName, err1 := getAppConfigName(name, active)
	if err1 != nil {
		log.Fatal(err1)
	}
	content, err := GetConfigByDataId(appConfigName)
	if err != nil {
		Logger.Info("getConfig error dataId=" + appConfigName)
	}
	err = viper.MergeConfig(strings.NewReader(content))
	if err != nil {
		Logger.Error("Viper Failed to resolve configuration content=" + content)
	}
	// 监听配置
	err = ListenConfig(appConfigName, LocalNacosConfig.Group, func(namespace, group, dataId, data string) {
		fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		err = viper.MergeConfig(strings.NewReader(data))
		if err != nil {
			Logger.Error("Viper Failed to resolve configuration content=" + content)
			return
		}

	})

	RegisterClient(name, LocalNacosConfig.ClusterName, LocalNacosConfig.Group)
	services, err := NamingClient.GetService(vo.GetServiceParam{
		ServiceName: name,
		Clusters:    []string{LocalNacosConfig.ClusterName}, // 默认值DEFAULT
		GroupName:   LocalNacosConfig.Group,                 // 默认值DEFAULT_GROUP
	})
	Logger.Debug(fmt.Sprintf("server=%s,instances=%v", name, services))
	if err != nil {
		fmt.Println(err)
		Logger.Info("ListenConfigByDataId failed")
		return
	}
}
func RegisterClient(serverName, cluster, group string) {
	if serverName != "" {
		ip, _ := GetOutBoundIP()
		// 注册自身到nacos
		success, err1 := NamingClient.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          ip,
			Port:        uint64(viper.GetInt("server.port")),
			ServiceName: serverName,
			Weight:      10,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			Metadata:    map[string]string{"idc": "shanghai"},
			ClusterName: cluster, // 默认值DEFAULT
			GroupName:   group,   // 默认值DEFAULT_GROUP
		})
		if err1 != nil {
			Logger.Error("nacos register failed")
		}
		if success {
			Logger.Info("nacos register success,should subscribe other server instance")
		}
	} else {
		Logger.Info("startNamingClient server name is empty")
	}

}

func GetConfigByDataId(dataId string) (content string, err error) {

	return GetConfig(dataId, DefaultGroup)
}

func GetConfig(dataId, group string) (content string, err error) {

	content, err = ConfigClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	return
}

type ConfigListenHandler func(namespace, group, dataId, data string)

func ListenConfig(dataId, group string, handler ConfigListenHandler) error {
	return ConfigClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			handler(namespace, group, dataId, data)
		},
	})
}

func NacosDestroy() {
	Logger.Info("nacos Destroy")
}
