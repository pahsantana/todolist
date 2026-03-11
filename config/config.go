package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	MongoURI    string
	MongoDB     string
	ServerPort  string
	Environment string
}

func Load() *Config {
	username := requireEnv("MONGO_ROOT_USERNAME")
	password := requireEnv("MONGO_ROOT_PASSWORD")
	host := requireEnv("MONGO_HOST")

	return &Config{
		MongoURI:    fmt.Sprintf("mongodb://%s:%s@%s:27017/", username, password, host),
		MongoDB:     requireEnv("MONGO_DATABASE"),
		ServerPort:  requireEnv("SERVER_PORT"),
		Environment: requireEnv("ENVIRONMENT"),
	}
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("required environment variable not set: %s", key)
	}
	return value
}
