package config

type Config struct {
	App AppConfig `config:"app"`
	Server Server `config:"server"`
	Database Database `config:"database"`
}

type Server struct {
	Host   string `yaml:"host"`
	Port int `yaml:"port"`
}

type Database struct {
	ConnPath string `yaml:"conn_path"`
}

type AppConfig struct {
	DebugStatus bool `yaml:"debug_status"`
}
