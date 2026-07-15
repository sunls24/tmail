package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host         string   `env:"HOST"`
	Port         string   `env:"PORT" envDefault:"3000"`
	DomainList   []string `env:"DOMAIN_LIST"`
	AdminAddress string   `env:"ADMIN_ADDRESS"`
	BaseDir      string   `env:"BASE_DIR" envDefault:"fs"`
	DB           Database `envPrefix:"DB_"`
	Debug        bool     `env:"DEBUG"`

	TurnstileSiteKey   string        `env:"TURNSTILE_SITE_KEY"`
	TurnstileSecretKey string        `env:"TURNSTILE_SECRET_KEY"`
	TurnstileCookieTTL time.Duration `env:"TURNSTILE_COOKIE_TTL" envDefault:"6h"`
}

type Database struct {
	Driver string `env:"DRIVER" envDefault:"postgres"`
	Host   string `env:"HOST"`
	Port   string `env:"PORT" envDefault:"5432"`
	User   string `env:"USER" envDefault:"postgres"`
	Pass   string `env:"PASS"`
	Name   string `env:"NAME" envDefault:"tmail"`
}

func MustNew() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	if (cfg.TurnstileSiteKey == "") != (cfg.TurnstileSecretKey == "") {
		panic("TURNSTILE_SITE_KEY and TURNSTILE_SECRET_KEY must be set together")
	}
	if cfg.TurnstileEnabled() && cfg.TurnstileCookieTTL <= 0 {
		panic(fmt.Sprintf("invalid TURNSTILE_COOKIE_TTL: %s", cfg.TurnstileCookieTTL))
	}
	return &cfg
}

func (c *Config) TurnstileEnabled() bool {
	return c.TurnstileSiteKey != "" && c.TurnstileSecretKey != ""
}
