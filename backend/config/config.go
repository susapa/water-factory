package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type JWTConfig struct {
	Secret     string
	AccessTTL  string
	RefreshTTL string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	cfg := &Config{
		App: AppConfig{
			Port: viper.GetString("APP_PORT"),
			Env:  viper.GetString("APP_ENV"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			Name:     viper.GetString("DB_NAME"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			AccessTTL:  viper.GetString("JWT_ACCESS_TTL"),
			RefreshTTL: viper.GetString("JWT_REFRESH_TTL"),
		},
	}

	if cfg.App.Port == "" {
		cfg.App.Port = "8080"
	}
	if cfg.App.Env == "" {
		cfg.App.Env = "development"
	}
	if cfg.JWT.AccessTTL == "" {
		cfg.JWT.AccessTTL = "15m"
	}
	if cfg.JWT.RefreshTTL == "" {
		cfg.JWT.RefreshTTL = "168h"
	}

	return cfg, nil
}
