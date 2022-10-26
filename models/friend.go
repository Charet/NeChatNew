package models

import (
	"NechatService/dao"
	"log"
)

type Friend struct {
	UserID   int `json:"UserID" db:"user_id"`
	FriendID int `json:"FriendID" db:"friend_id"`
}

type ApplyFriend struct {
	SenderID   int  `json:"SenderID" db:"sender_id"`
	ReceiverID int  `json:"ReceiverID" db:"receiver_id"`
	IsRead     bool `json:"IsRead" db:"is_read"`
}

// SaveApplyFriend 保存添加好友申请(发送人ID,接收人ID,申请是否已读)至数据库中
func SaveApplyFriend(applyFriend *ApplyFriend) error {
	sqlStr := "INSERT INTO apply_friend (sender_id, receiver_id, is_read) VALUES (?, ?, ?)"
	_, err := dao.DB.Exec(sqlStr, applyFriend.SenderID, applyFriend.ReceiverID, 0)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ChangeApplyFriendStatus 修改好友申请的状态为已读
func ChangeApplyFriendStatus(applyFriend *ApplyFriend) error {
	sqlStr := "UPDATE apply_friend SET is_read = 1 WHERE sender_id = ? AND receiver_id = ?"
	ret, err := dao.DB.Exec(sqlStr, applyFriend.SenderID, applyFriend.ReceiverID)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = ret.RowsAffected() // 此处垃圾桶变量存储的为操作影响的行数(后期若有log可以打印)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// HaveApplyFriend 查询好友请求中是否有userID
func HaveApplyFriend(userID int) (bool, error) {
	applyFriend := ApplyFriend{}
	sqlStr := "SELECT * FROM apply_friend WHERE receiver_id = ? AND is_read = 0 LIMIT 1"
	err := dao.DB.Get(&applyFriend, sqlStr, userID)
	if (err != nil) && (err.Error() != "sql: no rows in result set") {
		log.Println(err)
		return false, err
	} else if (err != nil) && (err.Error() == "sql: no rows in result set") {
		return false, nil
	}
	return true, nil
}

// GetApplyFriend 返回userID的所有好友申请
func GetApplyFriend(userID int) ([]ApplyFriend, error) {
	sqlStr := "SELECT sender_id FROM apply_friend WHERE receiver_id = ? AND is_read = 0"
	var applyFriend []ApplyFriend
	err := dao.DB.Select(&applyFriend, sqlStr, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return applyFriend, nil
}

// SaveFriend 保存好友信息
func SaveFriend(friend *Friend) error {
	sqlStr := "INSERT INTO friend (user_id, friend_id) VALUES (?, ?)"
	_, err := dao.DB.Exec(sqlStr, friend.UserID, friend.FriendID)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = dao.DB.Exec(sqlStr, friend.FriendID, friend.UserID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DeleteFriend 删除好友
func DeleteFriend(friend *Friend) error {
	sqlStr := "DELETE FROM friend WHERE user_id = ? AND friend_id = ?"
	_, err := dao.DB.Exec(sqlStr, friend.UserID, friend.FriendID)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = dao.DB.Exec(sqlStr, friend.FriendID, friend.UserID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetFriendList 查询好友列表
func GetFriendList(userID int) ([]Friend, error) {
	sqlStr := "SELECT * FROM friend WHERE user_id = ?"
	var friend []Friend
	err := dao.DB.Select(&friend, sqlStr, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return friend, err
}
