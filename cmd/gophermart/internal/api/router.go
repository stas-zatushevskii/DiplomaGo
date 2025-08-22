package api

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/handlers"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/middlewares"
	"go.uber.org/zap"
	"sync"
)

// New(metricService *service.MetricsService) FIXME

func NewRouter(logger *zap.Logger, db *sql.DB, wg *sync.WaitGroup) *chi.Mux {
	router := chi.NewRouter()
	handler := handlers.NewHandler(logger, db)
	router.Use(middlewares.WithWaitGroup(wg))

	router.Get("/", handler.TestHandler())

	return router
}
