package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const (
	PageSize    = 3
	MaxContacts = 5
)

var ExpirationTime = time.Duration(0)

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

type Redis struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Password   string `yaml:"password"`
	Expiration string `yaml:"expiration"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Logger   Logger   `yaml:"logger"`
	Redis    Redis    `yaml:"redis"`
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

	ExpirationTime, err = time.ParseDuration(cfg.Redis.Expiration)
	if err != nil {
		return nil, fmt.Errorf("парсинг redis expiration: %w", err)
	}

	return cfg, nil
}
