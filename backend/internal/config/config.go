package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	PageSize    = 3
	MaxContacts = 5
)

type Server struct {
	JwtKey      string `yaml:"jwt_key"`
	ServerHost  string `yaml:"server_host"`
	ServerPort  string `yaml:"server_port"`
	MetricsHost string `yaml:"metrics_host"`
	MetricsPort string `yaml:"metrics_port"`
}

type Database struct {
	Name     string `yaml:"db_name"`
	User     string `yaml:"db_user"`
	Password string `yaml:"db_password"`
	Driver   string `yaml:"db_driver"`
	Host     string `yaml:"db_host"`
	Port     string `yaml:"db_port"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Logger   Logger   `yaml:"logger"`
}

func ReadConfig() (cfg *Config, err error) {
	cfg = new(Config)

	var f *os.File
	f, err = os.Open("config.yml")
	if err != nil {
		return nil, fmt.Errorf("открытие файла конфига: %w", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("чтение файла конфига: %w", err)
	}

	return cfg, nil
}
