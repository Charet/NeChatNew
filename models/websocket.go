package models

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	WriteBufferSize: 10240, //byte
	ReadBufferSize:  10240,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Clients = make(map[string]*Client) //用户连接池

type Client struct { //用户连接
	Client    *websocket.Conn
	Broadcast chan *ReceiveMsgType
	Status    chan bool
}

type ReceiveMsgType struct { //接收信息
	From    string `json:"from"`    //发送者
	To      string `json:"to"`      //接收者
	Content string `json:"content"` //消息内容
}

// ProcessMsg 获取广播中的消息并发送给对应连接
func (c *Client) ProcessMsg() {
	defer log.Println("[websocket.go/ProcessMsg]: Goroutine Closed.")
	for data := range c.Broadcast {
		toClient := Clients[data.To]
		err := toClient.Client.WriteJSON(gin.H{"from": data.From, "content": data.Content})
		// TODO 聊天内容存储到数据库中
		if err != nil {
			log.Println("[Models/websocket.go/ProcessMsg/WriteJSON]: ", err)
		}
	}
}

func (c *Client) ReadIndMsg(userID string) {
	defer log.Println("[websocket.go/ReadIndMsg]: Goroutine Closed.")
	for {
		log.Println("[Models/websocket.go/ReadIndMsg]: User Connected, Reading Message.")
		_, msg, err := c.Client.ReadMessage()
		// websocket退出
		if err != nil {
			log.Println("[Models/websocket.go/ReadIndMsg/ReadMessage]: ", err)
			delete(Clients, userID) // 从Clients池中注销用户.
			return
		} else if len(msg) == 0 {
			log.Println("[Models/websocket.go/ReadIndMsg]: Don't have message.")
			continue
		}
		var receiveMsg ReceiveMsgType
		err = json.Unmarshal(msg, &receiveMsg)
		c.Broadcast <- &receiveMsg
	}
}
