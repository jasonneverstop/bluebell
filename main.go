package main

import (
	"bluebell_renjiexuan/controller"
	"bluebell_renjiexuan/dao/mysql"
	"bluebell_renjiexuan/dao/redis"
	"bluebell_renjiexuan/logger"
	"bluebell_renjiexuan/router"
	"bluebell_renjiexuan/setting"
	"fmt"
)

func main() {
	// 加载配置
	if err := setting.Init(); err != nil {
		fmt.Printf("load config failed,err:%v\n", err)
		return
	}
	// 初始化日志
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("logger.Init failed,err:%v\n", err)
		return
	}
	// 初始化Mysql
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		fmt.Printf("mysql.Init failed,err:%v\n", err)
		return
	}
	defer mysql.Close() // 程序退出关闭数据库连接
	// 初始化redis
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed,err:%v\n", err)
		return
	}
	defer redis.Close()
	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("Init validator failed,err:%v\n", err)
		return
	}
	// 注册路由
	r := router.SetupRouter(setting.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed,err:%v\n", err)
		retrun
	}

}
