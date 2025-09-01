package main

import (
	"context"
	"errors"
	"fmt"
	cfg "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api"
	db "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/database"
	CustomErrors "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
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
	logger.Info(fmt.Sprintf("%+v", config))
	if err != nil {
		logger.Error("failed to load config file", zap.Error(err))
		if errors.Is(err, CustomErrors.ErrConfigNotFound) {
			logger.Warn(fmt.Sprintf("App running on default settings %v", config))
		} else {
			return
		}
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
