package config

import (
	"encoding/json"
	"os"
)

var _conf AppConfig

type AppConfig struct {
	ServerPort int `json:"server_port"`
}

func Init() error {
	data, err := os.ReadFile("config/config.json")
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
