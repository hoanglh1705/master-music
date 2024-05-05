package config

import (
	"github.com/caarlos0/env/v5"
	"github.com/joho/godotenv"
)

// Configuration holds data necessary for configuring application
type Configuration struct {
	Stage        string   `env:"STAGE"`
	AppName      string   `env:"APP_NAME"`
	Port         int      `env:"PORT"`
	AllowOrigins []string `env:"ALLOW_ORIGINS"`
	Debug        bool     `env:"DEBUG"`
	DBUrl        string   `env:"DB_URL"`
	DBName       string   `env:"DB_NAME"`
}

// Load returns Configuration struct
func Load() (*Configuration, error) {
	cfg := new(Configuration)
	if err := PreloadLocalENV(); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// PreloadLocalENV reads .env* files and sets the values to os ENV
func PreloadLocalENV() error {
	basePath := ""
	// // local config per stage
	// if stage != "" {
	// 	godotenv.Load(basePath + ".env." + stage + ".local")
	// }

	// local config
	godotenv.Load(basePath + ".env.local")

	// // per stage config
	// if stage != "" {
	// 	godotenv.Load(basePath + ".env." + stage)
	// }

	// default config
	return godotenv.Load(basePath + ".env")
}
