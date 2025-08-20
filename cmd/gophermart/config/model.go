package config

type Config struct {
	App AppConfig `config:"app"`
	Server Server `config:"server"`
}

type Server struct {
	Host   string `yaml:"host"`
	Port int `yaml:"port"`
}

type AppConfig struct {
	DebugStatus bool `yaml:"debug_status"`
}
