package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/satori/go.uuid"
)

type Config struct {
	HostID         string `json:"host_id"`
	RedisServer    string `json:"redis_server"`
	CommandTimeout int    `json:"command_timeout"`

	HTTPPort string `json:"http_port"`
}

var Conf = Config{
	RedisServer:    ":6379",
	HostID:         uuid.NewV4().String(),
	CommandTimeout: 5,
	HTTPPort:       ":4000",
}

func Load() error {
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Println("File 'config.json' not found. Creating one with default configuration ...")
		configFile, err = json.MarshalIndent(Conf, "", "   ")
		if err != nil {
			log.Fatal("Failed to write config file!")
		}
		return ioutil.WriteFile("./config.json", configFile, 0660)
	}

	return json.Unmarshal(configFile, &Conf)
}
