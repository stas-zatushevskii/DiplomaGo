package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/handlers"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/middlewares"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/service"
	"go.uber.org/zap"
	"sync"
)

type RouterData struct{}

func NewRouter(logger *zap.Logger, service *service.Service, wg *sync.WaitGroup, orderChan chan<- models.ProcessOderData) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares.WithWaitGroup(wg))
	r.Use(middlewares.HeaderSetMiddleware)
	r.Use(middlewares.GzipMiddleware)

	validate := validator.New()
	h := handlers.NewHandler(logger, service, validate, orderChan)

	r.Route("/api/user", func(r chi.Router) {
		// public
		r.Post("/login", h.Login())
		r.Post("/register", h.Register())

		// auth
		r.Group(func(auth chi.Router) {
			auth.Use(middlewares.JWTMiddleware(service))
			auth.With(middlewares.CheckHeaderMiddleware).Post("/orders", h.OrderCreate())
			auth.Get("/orders", h.OrdersGet())
			auth.Get("/balance", h.GetUserBalance())
			auth.Post("/balance/withdraw", h.WithdrawOrderAccrual())
			auth.Get("/withdrawals", h.GetWithdrawalsHistory())
		})
	})

	return r
}
