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

	var reqWaitGroup sync.WaitGroup // FIXME

	db, err := database.NewDatabase(logger, cfg, &reqWaitGroup) // FIXME
	if err != nil {
		logger.Fatal("failed to create database", zap.Error(err))
		return
	}
	// TODO: add service
	router := api.NewRouter(logger, db.Db, &reqWaitGroup) // FIXME
	server := api.NewServer(router, logger, cfg)
	if err := server.Start(); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
		return
	}
	logger.Info("Server started")
	<-ctx.Done()
	GracefulShutdown(logger, server, db)
}

/*
ShutDown logic:
	listen for ctx.Done(), if got signal: creates new chan "done",
	run goroutine db.Close(done).

	goroutine db.Close(done):
		waiting wg.Done(), closing database, closing chan "done"
		(wg.Done() will happened when all active requests finish their job)

	when chan "done" is closed - exiting from main function
*/

func GracefulShutdown(logger *zap.Logger, server *api.Server, database *database.Database) {
	logger.Warn("shutdown: start")
	server.ServerShutdown()
	logger.Info("shutdown: server closed")
	database.DatabaseShutdown()
	logger.Info("shutdown: database closed")
	logger.Warn("shutdown: end")
}
