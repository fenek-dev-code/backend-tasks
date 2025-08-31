package main

import (
	"database/sql"
	"log"

	"github.com/fenek-dev/backend-tasks/internal/config"
	"github.com/fenek-dev/backend-tasks/internal/logger"
	"github.com/fenek-dev/backend-tasks/internal/storage"
	"go.uber.org/zap"
)

func main() {
	// TODO: init config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed to init or read config file: ", err)
	}

	// TODO: init logger
	loger, err := logger.NewLogger(cfg.LoggerConfig.LogLevl, cfg.LoggerConfig.LogPath)
	if err != nil {
		log.Fatal("failed to init logger: ", err)
	}
	loger.Info("logger: logger initialized")
	loger.Debug("logger: debug message")

	// TODO: init database
	db_url := cfg.DataBaseConfig.GetPostgresURL()
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		loger.Fatal("failed to init database: ", zap.Error(err))
	}

	// TODO: init storage
	storage := storage.NewStorage(db, loger)
	if err := storage.Ping(); err != nil {
		loger.Fatal("failed to ping database: ", zap.Error(err))
	}
	loger.Info("storage: storage initialized")

	// TODO: init server and router

	// TODO: start server
}
