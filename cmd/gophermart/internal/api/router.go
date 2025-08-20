package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/handlers"
	"go.uber.org/zap"
)

// New(metricService *service.MetricsService) FIXME

func NewRouter(logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	handler := handlers.NewHandler(logger)

	router.Get("/", handler.TestHandler())
	return router
}
