package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"go.uber.org/zap"
)

type Database struct {
	Db     *sql.DB
	logger *zap.Logger
}

func NewDatabase(logger *zap.Logger, config *config.Config) (*Database, error) {
	db, err := sql.Open("pgx", config.Database.ConnPath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Database{Db: db, logger: logger}, nil
}

func (d *Database) DatabaseShutdown() {
	err := d.Db.Close()
	if err != nil {
		d.logger.Fatal("failed to close database", zap.Error(err))
	}
}
