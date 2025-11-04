package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

type EnvConfig struct {
	LogLevel string `env:"LOG_LEVEL,required"`

	// Config
	Mode        string   `env:"GIN_MODE,required" envDefault:"debug"`
	CorsOrigins []string `env:"CORS_ORIGINS,required"`
	AppURL      string   `env:"APP_URL,required"`
	DemoMode    bool     `env:"DEMO_MODE" envDefault:"false"`

	// Http Port
	Port string `env:"HTTP_PORT,required"`

	// Database
	DSN string `env:"DB_DSN,required"`

	// Redis
	// RedisAddr     string `env:"REDIS_HOST,required"`
	// RedisDB       int    `env:"REDIS_DB,required"`
	// RedisPassword string `env:"REDIS_PASSWORD,required"`

	// Token settings
	TokenSecret             string        `env:"TOKEN_SECRET,required"`
	TokenAccessTTL          time.Duration `env:"TOKEN_ACCESS_TTL,required" envDefault:"15m"`
	TokenRefreshTTL         time.Duration `env:"TOKEN_REFRESH_TTL,required" envDefault:"24h"`
	TokenRefreshRememberTTL time.Duration `env:"TOKEN_REFRESH_REMEMBER_TTL,required" envDefault:"720h"`
}

func Load() (*EnvConfig, error) {
	var cfg EnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load env: %v", err)
	}

	return &cfg, nil
}
