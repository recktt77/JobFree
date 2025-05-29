package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	MongoURI  string
	MongoDB   string
	RedisAddr string
	RedisPwd  string
	NatsURL   string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port:      getEnv("PORT", "50052"),
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:   getEnv("MONGO_DB", "freelance_db"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPwd:  getEnv("REDIS_PASSWORD", ""),
		NatsURL:   getEnv("NATS_URL", "nats://localhost:4222"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
