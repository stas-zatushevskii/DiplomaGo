package config

import (
	"flag"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

type NewConfig struct {
	log    *zap.Logger
	config *Config
}

func (cfg *NewConfig) ParseEnv() error {
	// TODO
	return nil
} // FIXME: maybe unnecessary

func (cfg *NewConfig) ParseFlags() {

	flag.Parse()
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
	err = cfg.ParseEnv()
	cfg.ParseFlags()
	return cfg.config, err
}
