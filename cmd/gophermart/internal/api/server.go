package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	config *config.Config
	logger *zap.Logger
	router *chi.Mux
	ctx    context.Context
}

func NewServer(ctx context.Context, router *chi.Mux, logger *zap.Logger, config *config.Config) *Server {
	return &Server{
		logger: logger,
		ctx:    ctx,
		router: router,
		config: config,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Running server", zap.String("address", s.config.Server.Host+":"+strconv.Itoa(s.config.Server.Port)))
	srv := &http.Server{
		Addr:    s.config.Server.Host + ":" + strconv.Itoa(s.config.Server.Port),
		Handler: logger.WithLogging(s.router, s.logger),
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()
	select {
	case <-s.ctx.Done():
		s.logger.Warn("Shutting down server...")
		httpCtx, httpCancel := context.WithTimeout(context.Background(), 30*time.Second) // 30 sec for close all active requests
		defer httpCancel()
		return srv.Shutdown(httpCtx)
	case err := <-errCh:
		return err
	}
}
