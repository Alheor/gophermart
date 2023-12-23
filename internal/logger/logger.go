package logger

import (
	"github.com/Alheor/gophermart/internal/config"
	"go.uber.org/zap"
)

var logger *zap.Logger

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

	logger = zl
}

func GetLogger() *zap.Logger {
	return logger
}
