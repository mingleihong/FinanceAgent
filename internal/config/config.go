package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type LLMConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

type AppConfig struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	TTL    time.Duration `mapstructure:"ttl"`
}

type PostgresConfig struct {
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	DB              string        `mapstructure:"db"`
	ConnectTimeout  time.Duration `mapstructure:"connect_timeout"`
	MaxConns        int32         `mapstructure:"max_conns"`
	MinConns        int32         `mapstructure:"min_conns"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

type Config struct {
	JWT      *JWTConfig      `mapstructure:"jwt"`
	LLM      *LLMConfig      `mapstructure:"llm"`
	App      *AppConfig      `mapstructure:"app"`
	Postgres *PostgresConfig `mapstructure:"postgres"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config.yaml: %w", err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("decode config.yaml: %w", err)
	}
	if cfg.JWT == nil || len(cfg.JWT.Secret) < 32 {
		return nil, errors.New("jwt.secret must contain at least 32 bytes; set JWT_SECRET")
	}
	if cfg.JWT.TTL < time.Second {
		return nil, errors.New("jwt.ttl must be at least 1s")
	}
	return cfg, nil
}
