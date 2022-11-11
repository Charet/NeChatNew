package dao

import (
	"NechatService/config"
	"github.com/gomodule/redigo/redis"
	"log"
)

var RCoon redis.Conn

func InitRedis() {
	var err error
	RCoon, err = redis.Dial("tcp", config.ServerConfig.Redis.Host+":"+config.ServerConfig.Redis.Port)
	if err != nil {
		log.Println("[ERROR]Connect to redis failed,", err)
	}
}
