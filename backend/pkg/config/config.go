package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	Redis    RedisConfig    `mapstructure:",squash"`
}

type ServerConfig struct {
	Port string `mapstructure:"SERVER_PORT"`
	Mode string `mapstructure:"SERVER_MODE"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DBName   string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSL_MODE"`
	MaxConns int32  `mapstructure:"DB_MAX_CONNS"`
	MinConns int32  `mapstructure:"DB_MIN_CONNS"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("Error when read .env file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Error when unmarshal .env file: %w", err)
	}

	if cfg.Database.Host == "" || cfg.Database.Port == "" || cfg.Database.User == "" || cfg.Database.Password == "" || cfg.Database.DBName == "" || cfg.Database.SSLMode == "" || cfg.Database.MaxConns == 0 || cfg.Database.MinConns == 0 {
		return nil, fmt.Errorf("Missing required database configuration")
	}

	if cfg.Redis.Host == "" || cfg.Redis.Port == "" || cfg.Redis.Password == "" || cfg.Redis.DB == 0 {
		return nil, fmt.Errorf("Missing required redis configuration")
	}

	if cfg.Server.Port == "" || cfg.Server.Mode == "" {
		return nil, fmt.Errorf("Missing required server configuration")
	}

	log.Printf("Loading env config success %v", &cfg)

	return &cfg, nil
}