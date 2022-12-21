package main

import (
	"fmt"
	"os"

	config "github.com/lixvyang/Go-redis/configs"
	"github.com/lixvyang/Go-redis/internal/logger"
	RedisServer "github.com/lixvyang/Go-redis/internal/redis/server"
	"github.com/lixvyang/Go-redis/internal/tcp"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "github.com/lixvyang/Go-redis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	configFileName := os.Getenv("CONFIG")
	if configFileName == "" {
		if fileExists("configs/redis.conf") {
			config.SetupConfig("configs/redis.conf")
		}
	} else {
		config.SetupConfig(configFileName)
	}
	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, RedisServer.MakeHandler())
	if err != nil {
		logger.Error(err)
	}
}
