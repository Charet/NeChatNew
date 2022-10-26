package controller

import (
	"NechatService/logic"
	"NechatService/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func RegisterHandler(c *gin.Context) {
	userInfo := models.Userinfo{}
	jsonUnmarshal(c, &userInfo)

	c.JSON(logic.Register(&userInfo))
}

func LoginHandler(c *gin.Context) {
	var err error
	userInfo := models.Userinfo{}
	userInfo.UserID, err = strconv.Atoi(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": "UserID type error."})
		return
	}
	userInfo.Password = c.Param("password")

	c.JSON(logic.Login(&userInfo))
}

func DeleteAccountHandler(c *gin.Context) {
	userID, ok := c.Get("UserID")
	if !ok {
		log.Println("[ERROR]Handler var get failed.")
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "Server Error"})
		return
	}

	err := models.DeleteAccount(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
	} else {
		c.JSON(http.StatusNoContent, gin.H{"code": 0, "msg": "Account was delete."})
	}
}

func AddFriendHandler(c *gin.Context) {
	applyFriend := models.ApplyFriend{}
	jsonUnmarshal(c, &applyFriend)
	SenderID, ok := c.Get("UserID")
	if !ok {
		fmt.Println("[ERROR]Handler var get failed.")
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "Server Error"})
		return
	}
	applyFriend.SenderID = SenderID.(int)

	err := models.SaveApplyFriend(&applyFriend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func GetApplyFriendListHandler(c *gin.Context) {
	userID, ok := c.Get("UserID")
	if !ok {
		fmt.Println("[ERROR]Handler var get failed.")
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "Server Error"})
		return
	}

	fReq, err := models.GetApplyFriend(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
	}
	c.JSON(http.StatusOK, gin.H{"ApplyFriendList": fReq})
}

func ChangeApplyFriendStatusHandler(c *gin.Context) {
	applyFriend := models.ApplyFriend{}
	jsonUnmarshal(c, &applyFriend)

	err := models.ChangeApplyFriendStatus(&applyFriend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success."})
}

func ApplyFriendHandler(c *gin.Context) {
	friend := models.Friend{}
	jsonUnmarshal(c, &friend)

	err := models.SaveFriend(&friend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
	}
	c.JSON(http.StatusCreated, gin.H{"code": 0, "msg": "Add Friend success."})
}

func DeleteFriendHandler(c *gin.Context) {
	friend := models.Friend{}
	jsonUnmarshal(c, &friend)

	err := models.DeleteFriend(&friend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"code": 0, "msg": "Delete friend success."})
}

func GetFriendListHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
		return
	}
	list, err := models.GetFriendList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "Success.", "date": list})
}

/* 学习泛型练手 TODO: 封装为方法*/
type allStruct interface {
	*models.Userinfo | *models.ApplyFriend | *models.Friend
}

func jsonUnmarshal[T allStruct](c *gin.Context, data T) {
	resp, err := c.GetRawData()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "Get json failed"})
		return
	}

	err = json.Unmarshal(resp, data)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "Json unmarshal failed."})
		return
	}
}
