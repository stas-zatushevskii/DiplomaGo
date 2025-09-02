package config

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"strings"
)

type NewConfig struct {
	log    *zap.Logger
	config *Config
}

func ConvertPostgresURLToDSN(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	user := u.User.Username()
	password, _ := u.User.Password()

	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "5432"
	}

	dbname := strings.TrimPrefix(u.Path, "/")
	if dbname == "" {
		dbname = "postgres"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	return dsn, nil
}

func (cfg *Config) parseFlags() error {
	var (
		connPathFlag   string
		accrualAddrFlg string
		serverAddrFlg  string
	)

	flag.StringVar(&serverAddrFlg, "a", "", "address to listen (e.g. 127.0.0.1:8080)")
	flag.StringVar(&connPathFlag, "d", "", "database connection URL (postgres://...)")
	flag.StringVar(&accrualAddrFlg, "r", "", "accrual service URL/address")
	flag.Parse()

	if serverAddrFlg != "" {
		cfg.Server.Address = serverAddrFlg
	}

	if accrualAddrFlg != "" {
		cfg.Accrual.Address = accrualAddrFlg
	}

	if connPathFlag != "" {
		dsn, err := ConvertPostgresURLToDSN(connPathFlag)
		if err != nil {
			return fmt.Errorf("invalid -d value: %w", err)
		}
		cfg.Database.ConnPath = connPathFlag
		cfg.Database.Dsn = dsn
	}
	return nil
}

func (cfg *Config) parseVirtualEnvironment() error {
	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		cfg.Server.Address = envRunAddr
	}
	if envConnPath := os.Getenv("DATABASE_URI"); envConnPath != "" {
		dsn, err := ConvertPostgresURLToDSN(envConnPath)
		if err != nil {
			return fmt.Errorf("invalid DATABASE_URI: %w", err)
		}
		cfg.Database.ConnPath = envConnPath
		cfg.Database.Dsn = dsn
	}
	if envAccrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualAddr != "" {
		cfg.Accrual.Address = envAccrualAddr
	}
	return nil
}

func DefaultConfigBuilder() *Config {
	return &Config{
		Database: Database{
			ConnPath: "postgres://postgres:123@localhost:5432/postgres",
			Dsn:      "host=localhost user=postgres password=123 dbname=postgres port=5432 sslmode=disable",
		},
		Accrual: Accrual{
			Address: "localhost:8081",
		},
		App: AppConfig{
			DebugStatus:     true,
			NumberOfWorkers: 10,
		},
		Server: Server{
			Address: "localhost:8080",
		},
		JWT: JWT{
			Secret: "secret",
			Expire: 3,
		},
	}
}

func (c *NewConfig) ParseConfigFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		c.log.Warn("config file not found, falling back to defaults", zap.String("path", path), zap.Error(err))
		c.config = DefaultConfigBuilder()
		return nil
	}
	if err := yaml.Unmarshal(data, c.config); err != nil {
		c.log.Error("failed to parse config file", zap.Error(err))
		return err
	}
	return nil
}

func LoadConfig(log *zap.Logger) (*Config, error) {
	nc := &NewConfig{log: log, config: DefaultConfigBuilder()}

	if err := nc.ParseConfigFromFile("../../config/cfg.yml"); err != nil {
		return nil, err
	}

	if err := nc.config.parseVirtualEnvironment(); err != nil {
		return nil, err
	}

	if err := nc.config.parseFlags(); err != nil {
		return nil, err
	}

	log.Info("dsn", zap.String("dsn", nc.config.Database.Dsn))
	log.Info("db_url", zap.String("url", nc.config.Database.ConnPath))
	log.Info("accrual_addr", zap.String("addr", nc.config.Accrual.Address))

	return nc.config, nil
}
