package handlers

import (
	"database/sql"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Handler struct {
	// metricService *service.MetricsService TODO
	logger *zap.Logger
	db     *sql.DB
}

func NewHandler(log *zap.Logger, db *sql.DB) *Handler {
	return &Handler{logger: log, db: db}
}

func (h *Handler) TestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		err := h.db.Ping()
		if err != nil {
			h.logger.Error("failed to ping database", zap.Error(err))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}
}
