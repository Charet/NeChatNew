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

func LogoutAccountHandler(c *gin.Context) {
	userID, ok := c.Get("UserID")
	if !ok {
		log.Println("[ERROR]Handler var get failed.")
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "Server Error"})
		return
	}

	err := models.LogoutToken(strconv.Itoa(userID.(int)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Logout Success."})
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

	yes, err := models.HaveUser(applyFriend.ReceiverID) //???????????????????????????
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
		return
	}
	if yes {
		err = models.SaveApplyFriend(&applyFriend)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "Don't have user."})
	}
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
	var err error
	friend := models.Friend{}
	friend.FriendID, err = strconv.Atoi(c.Param("friendId"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
		return
	}

	err = models.DeleteFriend(&friend)
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

func NewSocketClientHandler(c *gin.Context) {
	// ??????ws??????
	conn, err := models.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": 2, "msg": err}) //TODO ??????????????????
		return
	}

	//?????????,???Clients????????????
	userIDAny, ok := c.Get("UserID")
	userID := strconv.Itoa(userIDAny.(int))

	if !ok {
		err := conn.Close()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": err})
			return
		}
		log.Println("[controller.go/NewSocketClientHandler/Get]: Don't have userid.")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": "Don't have userid."})
		return
	}
	newClient := models.Client{Client: conn, Broadcast: make(chan *models.ReceiveMsgType)}
	models.Clients[userID] = &newClient

	go newClient.ReadIndMsg(userID) //????????????
	go newClient.ProcessMsg()       //????????????
}
