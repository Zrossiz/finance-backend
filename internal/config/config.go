package config

import (
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/helpers"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		ENV      string `yaml:"env"`
		Severity uint32 `yaml:"severity"`
	} `yaml:"app"`

	Server struct {
		Addr               string `yaml:"addr"`
		JWTAccessSecret    string `yaml:"jwt_access_secret"`
		JWTRefreshSecret   string `yaml:"jwt_refresh_secret"`
		JWTAccessLifetime  string `yaml:"jwt_access_lifetime"`
		JWTRefreshLifetime string `yaml:"jwt_refresh_lifetime"`
	} `yaml:"server"`

	Postgres struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DB       string `yaml:"dbname"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"postgres"`
}

func New() (*Config, error) {
	configPath, err := helpers.WithBasePath("config")
	if err != nil {
		return nil, fmt.Errorf("get config base path err: %w", err)
	}

	filename := fmt.Sprintf("%s/config.yml", configPath)

	var cfg Config
	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
