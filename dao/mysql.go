package dao

import (
	"NechatService/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitSQL() {
	var err error
	dataSourceName := config.ServerConfig.SQL.User + ":" + config.ServerConfig.SQL.Pass + "@tcp(" + config.ServerConfig.SQL.Host + ":" + config.ServerConfig.SQL.Port + ")/" + config.ServerConfig.SQL.Database + "?charset=utf8mb4&parseTime=True"
	DB, err = sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		fmt.Println("[ERROR]Try to connect failed,", err)
	}
	DB.SetConnMaxLifetime(-1)
}
