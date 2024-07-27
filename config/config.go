package config

import (
	"encoding/json"
	"os"
)

var _conf AppConfig

type AppConfig struct {
	ServerPort int       `json:"server_port"`
	Db         *DBConfig `json:"db"`
}
type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
}

func Init() error {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &_conf)
	if err != nil {
		return err
	}
	return nil
}

func GetServerPort() int {
	return _conf.ServerPort
}

func GetDBConf() *DBConfig {
	return _conf.Db
}
