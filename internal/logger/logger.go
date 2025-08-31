package logger

import (
	"errors"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logLevel, logPath string) (*zap.Logger, error) {
	var config zap.Config
	if logPath == "" {
		return nil, errors.New("log path is empty")
	}

	if !exists(logPath) {
		log.Println("log file does not exist, creating...")
		if err := createFile(logPath); err != nil {
			log.Printf("error creating log file: %s", err)
			return nil, err
		}
	}

	if logLevel == "info" {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Настройки вывода
	config.OutputPaths = []string{"stdout", logPath}
	config.ErrorOutputPaths = []string{"stderr"}

	// Настройка encoder
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "level"

	// Добавление кастомных полей
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.Fields(
			zap.String("service", "my-rest-api"),
			zap.String("version", "1.0.0"),
		),
	}

	logger, err := config.Build(options...)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func createFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return err
	}
	defer file.Close() // Важно закрывать файл

	fmt.Printf("Файл %s создан успешно\n", filename)
	return nil
}
