package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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
		Addr: s.config.Server.Host + ":" + strconv.Itoa(s.config.Server.Port),
		// Handler: logger.WithLogging(r), TODO
		Handler: s.router,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()
	select {
	case <-s.ctx.Done():
		s.logger.Info("Shutting down server...")
		return srv.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}
