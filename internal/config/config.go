package config

import (
	"os"
	"strconv"

	"dify-wp-sync/internal/logger"
)

// Config holds all configuration for the application loaded from environment variables.
type Config struct {
	// OAuth related configs
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Port         string

	// Redis
	RedisAddr string
	RedisDB   int
	RedisPwd  string

	// Dify
	DifyToken   string
	DifyBaseURL string
}

// LoadConfig loads configuration from environment variables and performs basic validation.
func LoadConfig() (*Config, error) {
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	cfg := &Config{
		ClientID:     os.Getenv("WPCOM_CLIENT_ID"),
		ClientSecret: os.Getenv("WPCOM_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("WPCOM_REDIRECT_URI"),
		Port:         getEnv("PORT", "8080"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:      db,
		RedisPwd:     os.Getenv("REDIS_PASSWORD"),
		DifyToken:    os.Getenv("DIFY_API_KEY"),
		DifyBaseURL:  getEnv("DIFY_BASE_URL", "https://api.dify.ai/v1"),
	}

	// Validate critical fields
	if cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.RedirectURI == "" {
		logger.Log.Fatalf("Missing required OAuth configuration: CLIENT_ID, CLIENT_SECRET, and REDIRECT_URI must be set.")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
