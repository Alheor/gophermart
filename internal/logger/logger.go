package logger

import (
	"go.uber.org/zap"

	"github.com/Alheor/gophermart/internal/config"
)

type Logger struct {
	logger *zap.Logger
}

var log *Logger

func Init() {
	lvl, err := zap.ParseAtomicLevel(config.Options.LogLevel)
	if err != nil {
		panic(err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	defer zl.Sync()

	log = &Logger{logger: zl}
}

func (l *Logger) Panic(err error) {
	l.logger.Panic(err.Error())
}

func (l *Logger) Fatal(err error) {
	l.logger.Panic(err.Error())
}

func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Warn(msg)
}

func (l *Logger) Info(msg string) {
	l.logger.Warn(msg)
}

func GetLogger() *Logger {
	return log
}
