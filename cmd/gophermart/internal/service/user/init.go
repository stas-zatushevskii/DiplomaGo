package user

import (
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/database"
	"go.uber.org/zap"
)

type ServiceUser struct {
	config   *config.Config
	logger   *zap.Logger
	database *database.Database
}

func NewServiceUser(config *config.Config, logger *zap.Logger, database *database.Database) *ServiceUser {
	return &ServiceUser{
		config:   config,
		logger:   logger,
		database: database,
	}
}
