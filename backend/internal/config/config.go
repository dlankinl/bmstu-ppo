package config

import (
	"fmt"
	"os"
)

const PageSize = 20

type DBConfig struct {
	User     string
	Password string
	Database string
	Host     string
	Port     string
	Driver   string
}

type Config struct {
	JwtKey string
	DBConfig
}

func ReadConfig() (cfg *Config, err error) {
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		return nil, fmt.Errorf("JWT_KEY должен быть заполнен")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, fmt.Errorf("DB_USER должен быть заполнен")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD должен быть заполнен")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("DB_NAME должен быть заполнен")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("DB_HOST должен быть заполнен")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return nil, fmt.Errorf("DB_PORT должен быть заполнен")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		return nil, fmt.Errorf("DB_DRIVER должен быть заполнен")
	}

	dbCfg := DBConfig{
		User:     dbUser,
		Password: dbPassword,
		Database: dbName,
		Host:     dbHost,
		Port:     dbPort,
		Driver:   dbDriver,
	}

	return &Config{
		JwtKey:   jwtKey,
		DBConfig: dbCfg,
	}, nil
}
