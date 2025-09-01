package config

import (
	"flag"
	"fmt"
	CustomErrors "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
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

	flag.StringVar(&serverAddrFlg, "-a", "", "address to listen (e.g. 127.0.0.1:8080)")
	flag.StringVar(&connPathFlag, "-d", "", "database connection URL (postgres://...)")
	flag.StringVar(&accrualAddrFlg, "-r", "", "accrual service URL/address")
	flag.Parse()

	if serverAddrFlg != "" {
		cfg.Server.Address = serverAddrFlg
		fmt.Println("-----------------------------------------------------------")
		fmt.Println(serverAddrFlg)
		fmt.Println("-----------------------------------------------------------")
	}

	if accrualAddrFlg != "" {
		cfg.Accrual.Address = accrualAddrFlg
		fmt.Println("-----------------------------------------------------------")
		fmt.Println(accrualAddrFlg)
		fmt.Println("-----------------------------------------------------------")
	}

	if connPathFlag != "" {
		dsn, err := ConvertPostgresURLToDSN(connPathFlag)
		if err != nil {
			return fmt.Errorf("invalid -d value: %w", err)
		}
		cfg.Database.ConnPath = connPathFlag
		cfg.Database.Dsn = dsn
		fmt.Println("-----------------------------------------------------------")
		fmt.Println(dsn)
		fmt.Println("-----------------------------------------------------------")
	} else if cfg.Database.ConnPath != "" && cfg.Database.Dsn == "" {
		dsn, err := ConvertPostgresURLToDSN(cfg.Database.ConnPath)
		if err != nil {
			return fmt.Errorf("invalid ConnPath in file: %w", err)
		}
		cfg.Database.Dsn = dsn
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
			Address: "127.0.0.1:8080",
		},
		App: AppConfig{
			DebugStatus:     true,
			NumberOfWorkers: 10,
		},
		Server: Server{
			Address: "127.0.0.1:8080",
		},
		JWT: JWT{
			Secret: "secret",
		},
	}
}

func (c *NewConfig) ParseConfigFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		c.log.Warn("config file not found, falling back to defaults", zap.String("path", path), zap.Error(err))
		c.config = DefaultConfigBuilder()
		return CustomErrors.ErrConfigNotFound
	}
	if err := yaml.Unmarshal(data, c.config); err != nil {
		c.log.Error("failed to parse config file", zap.Error(err))
		return err
	}
	return nil
}

func LoadConfig(log *zap.Logger) (*Config, error) {
	nc := &NewConfig{log: log, config: new(Config)}

	_ = nc.ParseConfigFromFile("../../config/cfg.yml")

	if err := nc.config.parseFlags(); err != nil {
		return nil, err
	}
	return nc.config, nil
}
