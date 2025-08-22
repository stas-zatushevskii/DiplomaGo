package api

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	config        *config.Config
	logger        *zap.Logger
	router        *chi.Mux
	ServerCtx     context.Context
	ServerStopCtx context.CancelFunc
	Srv           *http.Server
}

func NewServer(router *chi.Mux, logger *zap.Logger, config *config.Config) *Server {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	return &Server{
		logger:        logger,
		ServerStopCtx: serverStopCtx,
		ServerCtx:     serverCtx,
		router:        router,
		config:        config,
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
	case err := <-errCh:
		return err
	}
}

func (s *Server) ServerShutdown() {
	s.logger.Info("Shutting down server...")
	shutdownCtx, _ := context.WithTimeout(s.ServerCtx, 30*time.Second)

	go func() {
		<-shutdownCtx.Done()
		if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
			s.logger.Fatal("graceful shutdown timed out.. forcing exit.")
		}
	}()

	// Trigger graceful shutdown
	err := s.Srv.Shutdown(shutdownCtx)
	if err != nil {
		s.logger.Error("error shutting down server", zap.Error(err))
	}
	s.ServerStopCtx()
}
