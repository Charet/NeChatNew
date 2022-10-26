package logic

import (
	"NechatService/config"
	"NechatService/models"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

func generateToken(userID int, username string) (signedString string, err error) {
	claim := models.MyClaim{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,               //签发日期(1分钟前
			ExpiresAt: time.Now().Unix() + 60*60*2,          //过期时间(2小时
			Issuer:    config.ServerConfig.Server.JWTIssuer, //签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedString, err = token.SignedString([]byte(config.ServerConfig.Server.JWTSingedKey))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return signedString, nil
}
