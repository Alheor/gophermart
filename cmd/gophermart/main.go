package main

import (
	"context"
	"net/http"

	"github.com/Alheor/gophermart/internal/accural"
	"github.com/Alheor/gophermart/internal/config"
	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/Alheor/gophermart/internal/router"
)

func main() {
	config.Load()
	logger.Init()

	ctx := context.Background()

	err := repository.Init(ctx)
	if err != nil {
		panic(err)
	}

	accural.Init(ctx)

	logger.GetLogger().Info("Starting server: " + config.Options.Addr)

	if config.Options.SignatureKey == config.DefaultLSignatureKey {
		logger.GetLogger().Warn("Used default signature key! Please change the key (-k option)!")
	}

	err = http.ListenAndServe(config.Options.Addr, router.Init())
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
}
