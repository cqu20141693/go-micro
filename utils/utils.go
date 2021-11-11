package utils

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net"
	"strings"
)

var Logger *zap.Logger

func init() {
	Logger, _ = zap.NewDevelopmentConfig().Build()
}
func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
func ToJSONString(o interface{}) string {

	marshal, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(marshal)
}