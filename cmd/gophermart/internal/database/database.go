package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/config"
	"go.uber.org/zap"
	"sync"
)

type Database struct {
	Db     *sql.DB
	logger *zap.Logger
	wg     *sync.WaitGroup
}

func NewDatabase(logger *zap.Logger, config *config.Config, wg *sync.WaitGroup) (*Database, error) {
	db, err := sql.Open("pgx", config.Database.ConnPath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Database{Db: db, logger: logger, wg: wg}, nil
}

func (d *Database) Close(done chan struct{}) {
	d.wg.Wait()
	d.logger.Warn("Closing database")

	err := d.Db.Close()
	if err != nil {
		d.logger.Fatal("failed to close database", zap.Error(err))
	}
	close(done)
}

func (d *Database) DatabaseShutdown() {
	err := d.Db.Close()
	if err != nil {
		d.logger.Fatal("failed to close database", zap.Error(err))
	}
}
