package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Db     *sql.DB
	GormDB *gorm.DB
	logger *zap.Logger
}

func NewDatabase(logger *zap.Logger, config *config.Config) (*Database, error) {
	db, err := sql.Open("pgx", config.Database.ConnPath)
	gormDB, err := gorm.Open(postgres.Open(config.Database.Dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	database := &Database{Db: db, GormDB: gormDB, logger: logger}
	return database, nil
}

func (d *Database) DatabaseShutdown() {
	err := d.Db.Close()
	if err != nil {
		d.logger.Fatal("failed to close database", zap.Error(err))
	}
}

func SetupDatabase(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Order{}, &models.OrderHistory{})
}
