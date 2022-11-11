package middleware

import (
	"NechatService/config"
	"NechatService/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"strconv"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "账户未登录"})
			c.Abort()
			return
		}
		t, err := jwt.ParseWithClaims(token, &models.MyClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.ServerConfig.Server.JWTSingedKey), nil
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "Login out"})
			c.Abort()
			return
		}
		trueToken, err := models.GetToken(strconv.Itoa(t.Claims.(*models.MyClaim).UserID))
		if err == redis.ErrNil { //数据库中没有token
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": "This Token can't use."})
			c.Abort()
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
			c.Abort()
			return
		}

		if trueToken == token { //接收到的token与数据库中的token的不符合
			c.Set("UserID", t.Claims.(*models.MyClaim).UserID)
			c.Set("Username", t.Claims.(*models.MyClaim).Username)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": "This Token can't use."})
			c.Abort()
			return
		}
	}
}
