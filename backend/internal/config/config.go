package config

import "os"

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	RedisAddr   string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/twitter_clone?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Port:        getEnv("PORT", "8080"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
