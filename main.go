package main

import (
	"encoding/json"
	"first-project/utils"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strings"
)

func main() {
	ip, err := utils.GetOutBoundIP()
	if err != nil {
		return
	}
	utils.Logger.Info(ip)
	content, err := utils.GetConfigByDataId("sip-link-dev.yml")
	if err != nil {
		utils.Logger.Info("getConfigByDataId error")
	}
	utils.Logger.Debug(content)
	err = utils.YamlViper.ReadConfig(strings.NewReader(content))
	if err != nil {
		utils.Logger.Error("Viper解析配置失败:")
	}
	stringMap := utils.YamlViper.GetStringMap("sip")

	marshal, _ := json.Marshal(stringMap)
	utils.Logger.Debug(string(marshal))
	update := make(chan bool)
	err = utils.ListenConfigByDataId("sip-link-dev.yml", func(namespace, group, dataId, data string) {
		fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		update <- true
	})
	serverName := utils.YamlViper.GetString("cc.application.name")
	startNamingClient(serverName, utils.Default_cluster, utils.Default_group)
	services, err := utils.NamingClient.GetService(vo.GetServiceParam{
		ServiceName: serverName,
		Clusters:    []string{utils.Default_cluster}, // 默认值DEFAULT
		GroupName:   utils.Default_group,             // 默认值DEFAULT_GROUP
	})
	sprintf := fmt.Sprintf("server=%s,instances=%v", serverName, services)
	utils.Logger.Debug(sprintf)
	if err != nil {
		fmt.Println(err)
		utils.Logger.Info("ListenConfigByDataId failed")
		return
	}
	if <-update {
		utils.Logger.Debug("update true")
	}
}

func startNamingClient(serverName, cluster, group string) {
	subscribeServerInstance(serverName, cluster, group)
	if serverName != "" {
		ip, _ := utils.GetOutBoundIP()
		success, err1 := utils.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          ip,
			Port:        8848,
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
			utils.Logger.Error("nacos register failed")
		}
		if success {
			utils.Logger.Info("nacos register success,should subscribe other server instance")
		}
	} else {
		utils.Logger.Info("startNamingClient server name is empty")
	}

}

func subscribeServerInstance(serverName string, cluster string, group string) {
	// Subscribe key=serviceName+groupName+cluster
	// 注意:我们可以在相同的key添加多个SubscribeCallback.
	utils.Err = utils.NamingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: serverName,
		GroupName:   group,             // 默认值DEFAULT_GROUP
		Clusters:    []string{cluster}, // 默认值DEFAULT
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			utils.Logger.Info(fmt.Sprintf("\n\n callback return services:%s \n\n", utils.ToJSONString(services)))
		},
	})
}
