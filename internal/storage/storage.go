package storage

import (
	"database/sql"

	"go.uber.org/zap"
)

type Storage struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewStorage(db *sql.DB, logger *zap.Logger) *Storage {
	logger.Info("storage: new storage initialized")
	return &Storage{
		DB:     db,
		Logger: logger,
	}
}

func (s *Storage) Ping() error {
	return s.DB.Ping()
}

func (s *Storage) Close() error {
	s.Logger.Info("storage: close storage")
	return s.DB.Close()
}
