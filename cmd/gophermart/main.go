package main

import (
	"context"
	cfg "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api"
	db "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/database"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
	srv "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/service"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/utils"
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

	config, err := cfg.LoadConfig(logger)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
		return
	}

	database, err := db.NewDatabase(logger, config)
	if err != nil {
		logger.Fatal("failed to create database", zap.Error(err))
		return
	}
	err = db.SetupDatabase(database.GormDB)
	if err != nil {
		logger.Fatal("failed to setup database", zap.Error(err))
		return
	}

	service := srv.NewService(config, logger, database)

	orderChan := make(chan models.ProcessOrderData, config.App.NumberOfWorkers)
	semaphore := utils.NewSemaphore(config.App.NumberOfWorkers, logger)
	go service.OrderService.OrderListener(ctx, orderChan, semaphore)
	go service.OrderService.OrderLoader(ctx, orderChan)

	var reqWg sync.WaitGroup
	router := api.NewRouter(logger, service, &reqWg, orderChan)
	server := api.NewServer(ctx, router, logger, config, &reqWg)
	server.Start()

	<-ctx.Done()
	StartGracefulShutdown(logger, server, database)
}

func StartGracefulShutdown(logger *zap.Logger, server *api.Server, database *db.Database) {
	logger.Warn("STARTED Graceful Shutdown")
	server.ServerShutdown()
	logger.Info("shutdown: server closed")
	database.DatabaseShutdown()
	logger.Info("shutdown: database closed")
	logger.Warn("ENDED Graceful Shutdown")
}
