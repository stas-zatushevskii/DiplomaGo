package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	// metricService *service.MetricsService TODO
	logger *zap.Logger
}

func NewHandler(log *zap.Logger) *Handler {
	return &Handler{logger: log}
}

func (h *Handler) TestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}
}