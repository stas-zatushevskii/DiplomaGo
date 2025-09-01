package config

type Config struct {
	App      AppConfig `config:"app"`
	Server   Server    `config:"server"`
	Database Database  `config:"database"`
	JWT      JWT       `config:"jwt"`
	Accrual  Accrual   `config:"accrual"`
}

type Server struct {
	Address string `yaml:"address"`
}

type Accrual struct {
	Address string `yaml:"address"`
}

type Database struct {
	ConnPath string `yaml:"conn_path"`
	Dsn      string `yaml:"dsn"`
}

type AppConfig struct {
	DebugStatus     bool `yaml:"debug_status"`
	NumberOfWorkers int  `yaml:"number_of_workers"`
}

type JWT struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}
