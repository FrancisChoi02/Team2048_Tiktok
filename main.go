package main

import (
	"Team2048_Tiktok/conf"
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/dao/redis"
	"Team2048_Tiktok/handler"
	"Team2048_Tiktok/logger"
	"Team2048_Tiktok/router"
	"fmt"
)

// @title Team2048_Tiktok 项目接口文档
// @version 1.0
// @description 极简版抖音

// @contact.name FrancisChoi
// @contact.url https://github.com/FrancisChoi02/
// @host 127.0.0.1:8080
// @BasePath
func main() {

	// 加载配置
	if err := conf.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	//加载日志器
	if err := logger.Init(conf.Conf.LogConfig, conf.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	//加载MySQL
	if err := mysql.Init(conf.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() // 程序退出关闭mysql连接

	//加载Redis
	if err := redis.Init(conf.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close() // 程序退出关闭redis数据库连接

	//加载validator参数验证器
	//验证用户用户名、密码的合法性
	if err := handler.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
		return
	}

	// 注册路由
	r := router.SetupRouter(conf.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", conf.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
