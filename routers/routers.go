package routers

import (
	"NechatService/controller"
	"NechatService/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.Cors())

	/*用户处理*/
	router.POST("/api/accounts", controller.RegisterHandler)                              //注册
	router.GET("/api/accounts/:userid/:password", controller.LoginHandler)                //登录
	router.DELETE("/api/accounts", middleware.JwtAuth(), controller.DeleteAccountHandler) //注销

	/*好友申请*/
	router.GET("/api/accounts/friend-requests", middleware.JwtAuth(), controller.GetApplyFriendListHandler)        //获取好友申请列表
	router.PATCH("/api/accounts/friend-requests", middleware.JwtAuth(), controller.ChangeApplyFriendStatusHandler) //修改好友申请状态为已读
	router.POST("/api/accounts/friend-requests", middleware.JwtAuth(), controller.ApplyFriendHandler)              //同意好友申请

	/*好友*/
	router.POST("/api/accounts/friends", middleware.JwtAuth(), controller.AddFriendHandler)            //添加好友
	router.DELETE("/api/accounts/friends", middleware.JwtAuth(), controller.DeleteFriendHandler)       //删除好友
	router.GET("/api/accounts/friends/:userid", middleware.JwtAuth(), controller.GetFriendListHandler) //获取好友列表

	router.GET("/api/ws", middleware.JwtAuth(), controller.NewSocketClientHandler)

	return router
}
