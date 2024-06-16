package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const PageSize = 3

type Config struct {
	Server struct {
		JwtKey string `yaml:"jwt_key"`
	} `yaml:"server"`
	Database struct {
		Name     string `yaml:"db_name"`
		User     string `yaml:"db_user"`
		Password string `yaml:"db_password"`
		Driver   string `yaml:"db_driver"`
		Host     string `yaml:"db_host"`
		Port     string `yaml:"db_port"`
	} `yaml:"database"`
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
