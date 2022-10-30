package models

import (
	"NechatService/dao"
	"log"
)

type Userinfo struct {
	UserID         int    `json:"UserID" db:"user_id"`
	Username       string `json:"Username" db:"username"`
	Password       string `json:"Password"`
	CryptoPassword string `db:"crypto_password"`
}

// SaveUserInfo 保存用户用户名以及加密后的密码
func SaveUserInfo(userInfo *Userinfo) error {
	sqlStr := "INSERT INTO user_info (username, crypto_password) VALUES (?, ?)"
	_, err := dao.DB.Exec(sqlStr, userInfo.Username, userInfo.CryptoPassword)
	if err != nil {
		log.Println("[ERROR]Inset data failed,", err)
		return err
	}
	return nil
}

// GetUserID 通过用户名以及加密后的密码查找用户ID
func GetUserID(userInfo *Userinfo) error {
	sqlStr := "SELECT user_id FROM user_info WHERE username = ? AND crypto_password = ?"
	err := dao.DB.Get(&userInfo.UserID, sqlStr, userInfo.Username, userInfo.CryptoPassword)
	if err != nil {
		log.Println("[models/user.go/GetUserID/Get]: ", err)
		return err
	}
	return nil
}

// GetUserInfoByID 通过用户ID查找用户名以及加密后的密码
func GetUserInfoByID(userInfo *Userinfo) error {
	sqlStr := "SELECT username, crypto_password FROM user_info WHERE user_id = ?"
	err := dao.DB.Get(userInfo, sqlStr, userInfo.UserID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteAccount(userID int) error {
	sqlStr := "DELETE FROM user_info WHERE user_id=?"
	_, err := dao.DB.Exec(sqlStr, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
