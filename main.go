package main

import (
	"bluebell_renjiexuan/setting"
	"fmt"
)

func main() {
	//加载配置
	if err := setting.Init(); err != nil {
		fmt.Printf("load config failed,err:%v\n", err)
		return
	}
}
