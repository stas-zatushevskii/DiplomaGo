package service

import (
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/database"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/service/order"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/service/user"
	"go.uber.org/zap"
)

type Service struct {
	UserService  *user.ServiceUser
	OrderService *order.ServiceOrder
}

func NewService(cfg *config.Config, logger *zap.Logger, database *database.Database) *Service {
	return &Service{
		UserService:  user.NewServiceUser(cfg, logger, database),
		OrderService: order.NewServiceOrder(cfg, logger, database),
	}
}
