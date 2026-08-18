package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	HTTPAddr           string
	DatabasePath       string
	AdminToken         string
	CORSOrigins        []string
	BotDriver          string
	WechatyPuppetToken string
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:           env("HTTP_ADDR", ":8080"),
		DatabasePath:       env("DATABASE_PATH", "./data/marc-chatbot.db"),
		AdminToken:         os.Getenv("ADMIN_TOKEN"),
		CORSOrigins:        splitCSV(env("CORS_ORIGINS", "http://localhost:5173")),
		BotDriver:          strings.ToLower(env("BOT_DRIVER", "mock")),
		WechatyPuppetToken: os.Getenv("WECHATY_PUPPET_SERVICE_TOKEN"),
	}
	if strings.TrimSpace(cfg.AdminToken) == "" {
		return Config{}, errors.New("ADMIN_TOKEN is required")
	}
	if cfg.BotDriver != "mock" && cfg.BotDriver != "wechaty" {
		return Config{}, errors.New("BOT_DRIVER must be mock or wechaty")
	}
	if cfg.BotDriver == "wechaty" && strings.TrimSpace(cfg.WechatyPuppetToken) == "" {
		return Config{}, errors.New("WECHATY_PUPPET_SERVICE_TOKEN is required for wechaty driver")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
