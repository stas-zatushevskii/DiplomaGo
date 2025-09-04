package order

import (
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/database"
	"go.uber.org/zap"
)

type ServiceOrder struct {
	config                *config.Config
	logger                *zap.Logger
	database              *database.Database
	ProcessingOrdersCache map[string]struct{}
}

func NewServiceOrder(config *config.Config, logger *zap.Logger, database *database.Database) *ServiceOrder {
	return &ServiceOrder{
		config:                config,
		logger:                logger,
		database:              database,
		ProcessingOrdersCache: make(map[string]struct{}),
	}
}
