package main

import (
	"os"
	"Go-redis/lib/logger"
	"Go-redis/config"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}


func main() {
	logger.Setup(&logger.Settings{
		Path: "logs",
		Name: "go-redis",
		Ext: "log",
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

	

}
