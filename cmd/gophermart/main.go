package main

import (
	"context"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/database"
	log "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := log.CreateLogger()

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
		return
	}

	db, err := database.NewDatabase(logger, cfg)
	if err != nil {
		logger.Fatal("failed to create database", zap.Error(err))
		return
	}

	// TODO: add service

	var reqWg sync.WaitGroup
	router := api.NewRouter(logger, db.Db, &reqWg)
	server := api.NewServer(ctx, router, logger, cfg, &reqWg)
	server.Start()

	<-ctx.Done()
	StartGracefulShutdown(logger, server, db)
}

func StartGracefulShutdown(logger *zap.Logger, server *api.Server, database *database.Database) {
	logger.Warn("STARTED Graceful Shutdown")
	server.ServerShutdown()
	logger.Info("shutdown: server closed")
	database.DatabaseShutdown()
	logger.Info("shutdown: database closed")
	logger.Warn("ENDED Graceful Shutdown")
}
