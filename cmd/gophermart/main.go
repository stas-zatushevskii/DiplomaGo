package main

import (
	"context"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api"
	log "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func main() {
	logger := log.CreateLogger()
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
		return
	}

	ctx, Done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer Done()
	// TODO: add database
	// TODO: add service
	router := api.NewRouter(logger)
	server := api.NewServer(ctx, router, logger, cfg)

	if err := server.Start(); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
	<-ctx.Done()
}
