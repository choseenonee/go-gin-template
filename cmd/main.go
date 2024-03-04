package main

import (
	"template/internal/delivery"
	"template/pkg/config"
	"template/pkg/database"
	"template/pkg/database/cached"
	"template/pkg/log"
)

func main() {
	// TODO: заполнить .env и .env.example

	logger, loggerInfoFile, loggerErrorFile := log.InitLogger()
	defer loggerInfoFile.Close()
	defer loggerErrorFile.Close()

	logger.Info("Logger Initialized")

	config.InitConfig()
	logger.Info("Config Initialized")

	db := database.GetDB()
	logger.Info("Database Initialized")

	redisSession := cached.InitRedis()
	logger.Info("Redis Initialized")

	//queue := rabbit.NewQueueHandler()
	//logger.Info("Rabbit queue Initialized")

	delivery.Start(
		db,
		logger,
		redisSession,
	)

}
