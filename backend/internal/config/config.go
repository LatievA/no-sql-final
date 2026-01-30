package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI       string
	DBName         string
	JWTSecret      string
	Port           string
	GinMode        string
	AllowedOrigins string
}

var AppConfig *Config

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		MongoURI:       getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:         getEnv("DB_NAME", "bicycle_store"),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-key"),
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
	}

	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
