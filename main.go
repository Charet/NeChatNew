package main

import (
	"NechatService/config"
	"NechatService/dao"
	"NechatService/routers"
	"fmt"
)

func main() {
	config.InitConfig()
	//gin.SetMode(gin.ReleaseMode) //release mode on
	dao.InitSQL()
	dao.InitRedis()

	router := routers.SetupRouter()
	err := router.Run(":" + config.ServerConfig.Server.Port)
	if err != nil {
		fmt.Println("[ERR] HTTP server start failed,", err)
		return
	}
}
