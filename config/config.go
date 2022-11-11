package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var ServerConfig = Config{}

type Config struct {
	Server Server `json:"Server"`
	SQL    SQL    `json:"SQL"`
	Redis  Redis  `json:"Redis"`
}

type Server struct {
	Port         string `json:"Port"`
	JWTSingedKey string `json:"JWTSingedKey"`
	JWTIssuer    string `json:"JWTIssuer"`
}

type SQL struct {
	User     string `json:"User"`
	Pass     string `json:"Pass"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	Database string `json:"Database"`
}
type Redis struct {
	Host string `json:"Host"`
	Port string `json:"Port"`
}

func InitConfig() {
	configFile, err := os.Open("config/config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(configFile)
	byteValue, _ := io.ReadAll(configFile)
	_ = json.Unmarshal(byteValue, &ServerConfig)
}
