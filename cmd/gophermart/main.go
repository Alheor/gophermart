package main

import (
	"context"
	"github.com/Alheor/gophermart/internal/accural"
	"github.com/Alheor/gophermart/internal/config"
	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/Alheor/gophermart/internal/router"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	config.Load()
	logger.Init()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repository.Init(ctx)
	if err != nil {
		panic(err)
	}

	accural.Init()

	logger.GetLogger().Info("Starting server", zap.String("addr", config.Options.Addr))

	if config.Options.SignatureKey == config.DefaultLSignatureKey {
		logger.GetLogger().Warn("Used default signature key! Please change the key (-k option)!")
	}

	err = http.ListenAndServe(config.Options.Addr, router.Init())
	if err != nil {
		logger.GetLogger().Fatal(err.Error())
	}
}
