package micro

import (
	"fmt"
	"github.com/cqu20141693/go-micro/container"
	"github.com/cqu20141693/go-micro/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	micro.LogInit()
	micro.NacosInit()

}

var exit = make(chan bool, 1)

type ConfigHandler func()

type BeanConfig func() (prefix string, o interface{})

var chs = make([]ConfigHandler, 8)

var bcs = make([]BeanConfig, 8)

func Run(args []string) {
	// 执行destroy
	defer Destroy()

	SetupSignalHandler(shutdown) // 注册监听信号，绑定信号处理机制
	micro.Logger.Info(fmt.Sprintf("start app args=%v", args))
	// 配置读取NacosInit方法内部实现

	// 数据绑定

	// 配置处理器
	for i := range chs {
		chs[i]()
	}
	for i := range bcs {
		container.InjectSingleton(bcs[i]())
	}
	// 启动web监听
	<-exit
}

func Destroy() {
	micro.NacosDestroy()
	container.Destroy()
}

func SetupSignalHandler(shutdownFunc func(bool)) {
	closeSignalChan := make(chan os.Signal, 1)
	// 监听四种关闭信号
	signal.Notify(closeSignalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-closeSignalChan
		log.Printf("got signal to exit [signal = %v]", sig)
		//判断关闭信号是否为SIGQUIT(用户发送Ctrl+/即可触发)
		shutdownFunc(sig == syscall.SIGQUIT)
	}()
}

func shutdown(isgraceful bool) {
	if isgraceful {
		micro.Logger.Info("graceful shutdown application")
		return
		//当满足 sig == syscall.SIGQUIT,做相应退出处理
	}
	// 不是syscall.SIGQUIT的退出信号时，做相应退出处理
	micro.Logger.Info("ungraceful shutdown application")
	exit <- true
}
