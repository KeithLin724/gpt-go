package pkg

import (
	"log/slog"
	"os"
)

type Log struct {
	infoLogger  *slog.Logger
	errorLogger *slog.Logger
}

func NewLog() *Log {
	infoOption := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	infoHandler := slog.NewTextHandler(os.Stdout, &infoOption)
	infoLogger := slog.New(infoHandler)
	errorOption := slog.HandlerOptions{
		Level: slog.LevelError,
	}
	errorHandler := slog.NewTextHandler(os.Stderr, &errorOption)
	errorLogger := slog.New(errorHandler)
	return &Log{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
	}
}

func (log *Log) Info(msg string) {
	log.infoLogger.Info(msg)
}

func (log *Log) Infof(msg string, args ...interface{}) {
	log.infoLogger.Info(msg, args...)
}

func (log *Log) Error(msg string) {
	log.errorLogger.Error(msg)
}

func (log *Log) Errorf(msg string, args ...interface{}) {
	log.errorLogger.Error(msg, args...)
}
