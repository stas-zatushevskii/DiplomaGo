package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/logger"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	config *config.Config
	logger *zap.Logger
	router *chi.Mux
	ctx    context.Context
	reqWg  *sync.WaitGroup
	Srv    *http.Server
}

func NewServer(ctx context.Context, router *chi.Mux, logger *zap.Logger, config *config.Config, wg *sync.WaitGroup) *Server {
	return &Server{
		logger: logger,
		ctx:    ctx,
		router: router,
		config: config,
		reqWg:  wg,
	}
}

func (s *Server) Start() {
	s.logger.Info("Starting server: ", zap.String("address", s.config.Server.Address))
	srv := &http.Server{
		Addr:    s.config.Server.Address,
		Handler: logger.WithLogging(s.router, s.logger),
	}
	s.Srv = srv
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Failed to start server", zap.Error(err))
			return
		}
	}()
	s.logger.Info("Server started")
}

func (s *Server) ServerShutdown() {
	shutdownCtx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	if err := s.Srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("error shutting down server", zap.Error(fmt.Errorf("server still processing old requests")))
	}

	s.reqWg.Wait()

}
