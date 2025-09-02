package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/fenek-dev-code/backend-tasks/internal/config"
	"github.com/fenek-dev-code/backend-tasks/internal/controller/router"
	"github.com/fenek-dev-code/backend-tasks/internal/logger"
	"github.com/fenek-dev-code/backend-tasks/internal/storage"

	_ "github.com/mattn/go-sqlite3"
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
	db, err := sql.Open("sqlite3", cfg.SqliteConfig.Path)
	if err != nil {
		loger.Fatal("failed to init database: ", zap.Error(err))
	}
	loger.Info("db: database initialized")

	// TODO: init storage
	storage := storage.NewStorage(db, loger)
	loger.Info("storage: storage initialized")

	// TODO: init server and router
	addr := cfg.ServerConfig.Host + ":" + fmt.Sprint(cfg.ServerConfig.Port)
	router := router.NewRouter(storage, addr, loger)
	router.Init()

	// TODO: start server
	if err := router.Run(); err != nil {
		loger.Fatal("failed to start server: ", zap.Error(err))
	}
	loger.Info("server: server started")

	storage.Close()
	loger.Info("server: server stopped")
}
