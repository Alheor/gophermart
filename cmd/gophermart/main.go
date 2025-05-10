package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alheor/gophermart/internal/accural"
	"github.com/Alheor/gophermart/internal/config"
	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/Alheor/gophermart/internal/router"
	"github.com/Alheor/gophermart/internal/shutdown"
	"github.com/Alheor/gophermart/internal/userauth"

	"go.uber.org/zap"
)

var shutdownTimeout = 5 * time.Second

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(``, err.(error))
			logger.Sync()
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdown.Init()
	cfg := config.Load()

	var err error

	userauth.Init(cfg.SignatureKey)
	err = logger.Init()
	if err != nil {
		panic(err)
	}

	if cfg.SignatureKey == config.DefaultLSignatureKey {
		logger.Error(`Used default signature key! Please change the key!`, nil)
	}

	if len(cfg.SignatureKey) == 0 {
		logger.Fatal(`Signature key is empty`, nil)
	}

	err = repository.Init(ctx, cfg.DatabaseURI)
	if err != nil {
		logger.Fatal(`repository init error`, nil)
	}

	accural.InitConnector(cfg.AccrualAddr)
	accural.InitService(ctx)

	srv := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: router.GetRoutes(),
	}

	shutdown.GetCloser().Add(srv.Shutdown)

	go func() {
		logger.Info("Starting server", zap.String("addr", cfg.RunAddr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(`error while starting http server`, err)
		}
	}()

	<-ctx.Done()

	println("shutting down ...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	shutdown.GetCloser().Close(shutdownCtx)
}
