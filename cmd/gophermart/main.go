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
	logger := log.CreateLogger()

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
		return
	}

	var reqWaitGroup sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	db, err := database.NewDatabase(logger, cfg, &reqWaitGroup) // TODO graceful shut  down for server and base
	if err != nil {
		logger.Fatal("failed to create database", zap.Error(err))
		return
	}
	// TODO: add service
	router := api.NewRouter(logger, db.Db, &reqWaitGroup)
	server := api.NewServer(ctx, router, logger, cfg)
	if err := server.Start(); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}

	<-ctx.Done()

	done := make(chan struct{})
	go db.Close(done)

	<-done

	logger.Info("shutdown: done")
}
