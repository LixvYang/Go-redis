package main

import (
	"Go-redis/config"
	"Go-redis/lib/logger"
	RedisServer "Go-redis/redis/server"
	"Go-redis/tcp"
	"fmt"
	"os"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "go-redis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	configFileName := os.Getenv("CONFIG")
	if configFileName == "" {
		if fileExists("redis.conf") {
			config.SetupConfig("redis.conf")
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
