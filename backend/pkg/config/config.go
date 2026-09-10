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

// DSN returns the PostgreSQL connection string.
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("error when reading .env file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error when unmarshaling config: %w", err)
	}

	// Validate
	if cfg.Database.Host == "" || cfg.Database.Port == "" || cfg.Database.User == "" ||
		cfg.Database.Password == "" || cfg.Database.DBName == "" || cfg.Database.SSLMode == "" ||
		cfg.Database.MaxConns <= 0 || cfg.Database.MinConns <= 0 {
		return nil, fmt.Errorf("missing required database configuration")
	}

	if cfg.Redis.Host == "" || cfg.Redis.Port == "" || cfg.Redis.DB < 0 {
		return nil, fmt.Errorf("missing required redis configuration")
	}

	if cfg.Server.Port == "" || cfg.Server.Mode == "" {
		return nil, fmt.Errorf("missing required server configuration")
	}

	log.Printf("Loading env config success: %+v", cfg)

	return &cfg, nil
}

func setDefaults() {
	// Server
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_MODE", "debug")

	// Database
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "gamemate")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "gamemate")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_MAX_CONNS", 20)
	viper.SetDefault("DB_MIN_CONNS", 5)

	// Redis
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
}