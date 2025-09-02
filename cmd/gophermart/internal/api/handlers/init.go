package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	logger    *zap.Logger
	service   *service.Service
	validator *validator.Validate
	orderChan chan<- string
}

func NewHandler(
	log *zap.Logger,
	service *service.Service,
	validator *validator.Validate,
	orderChan chan<- string) *Handler {
	return &Handler{logger: log, service: service, validator: validator, orderChan: orderChan}
}

func (h *Handler) Test() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
