package config

import "os"

type Config struct {
	MongoURI  string
	DBName    string
	JWTSecret string
	Port      string
}

func Load() *Config {
	return &Config{
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:    getEnv("DB_NAME", "twitter_clone"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
