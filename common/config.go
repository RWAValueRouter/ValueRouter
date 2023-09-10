package common

import (
	"encoding/json"
	"github.com/RWAValueRouter/FastMulThreshold-DSA/log"
	"os"
)

type config struct {
	DbConfig    DBConfig
	RedisConfig RedisConfig
}

type DBConfig struct {
	DbDriverName   string
	DbDriverSource string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
}

var (
	Conf       = new(config)
	Configfile string
)

func Init() {
	name := "config.json"
	if Configfile != "" {
		name = Configfile
	}
	if _, err := os.Stat(name); err != nil {
		log.Error("config file is not exist", "name", name)
		os.Exit(0)
	}
	file, _ := os.Open(name)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Conf)
	if err != nil {
		log.Error("Error init Config :", "msg", err.Error())
		os.Exit(0)
	}
}
