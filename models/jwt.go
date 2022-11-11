package models

import (
	"NechatService/dao"
	"github.com/golang-jwt/jwt"
	"github.com/gomodule/redigo/redis"
	"log"
)

type MyClaim struct {
	UserID   int
	Username string
	jwt.StandardClaims
}

func RegisterToken(userID string, token string) bool {
	//存入Key在seconds秒后删除：SETEX KEY SECONDS VALUES
	reply, err := dao.RCoon.Do("SETEX", userID, 60*60*2, token) //此处匿名变量存储语句返回
	if err != nil {
		log.Println(err)
		return false
	}
	if reply == "OK" {
		return true
	} else {
		return false
	}
}

func LogoutToken(userID string) error {
	_, err := dao.RCoon.Do("DEL", userID) //返回1为成功删除，0为找不到字段
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetToken(userID string) (string, error) {
	reply, err := redis.String(dao.RCoon.Do("GET", userID))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return reply, nil
}
