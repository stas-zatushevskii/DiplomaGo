package config

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"log"
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

func (cfg *Config) parseFlags() {
	var (
		connPath   string
		accrualAdr string
		serverAddr string
	)
	flag.StringVar(&serverAddr, "a", "127.0.0.1:8080", "address to listen")
	flag.StringVar(&connPath, "d", "", "database connection path")
	flag.StringVar(&accrualAdr, "r", "", "accrual url")
	flag.Parse()

	if connPath != "" {
		dsn, err := ConvertPostgresURLToDSN(cfg.Database.ConnPath)
		if err != nil {
			log.Fatal(err)
			return
		}
		cfg.Database.Dsn = dsn
	}
	if accrualAdr != "" {
		cfg.Accrual.Address = accrualAdr
	}
	if serverAddr != "127.0.0.1:8080" {
		cfg.Accrual.Address = serverAddr
	}
}

func (cfg *NewConfig) ParseConfig() error {
	f, err := os.ReadFile("../../config/cfg.yml")
	if err != nil {
		cfg.log.Error("failed to read config file", zap.Error(err))
		return err
	}
	if err := yaml.Unmarshal(f, &cfg.config); err != nil {
		cfg.log.Error("failed to parse config file", zap.Error(err))
		return err
	}
	return nil
}

func LoadConfig(log *zap.Logger) (*Config, error) {
	cfg := &NewConfig{
		log:    log,
		config: new(Config),
	}
	err := cfg.ParseConfig()
	cfg.config.parseFlags()
	return cfg.config, err
}
