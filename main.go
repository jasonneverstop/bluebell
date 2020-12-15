package main

import (
	"bluebell_renjiexuan/logger"
	"bluebell_renjiexuan/setting"
	"fmt"
)

func main() {
	//加载配置
	if err := setting.Init(); err != nil {
		fmt.Printf("load config failed,err:%v\n", err)
		return
	}
	//初始化日志
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("logger.Init failed,err:%v\n", err)
		return
	}
}
