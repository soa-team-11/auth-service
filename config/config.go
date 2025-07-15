package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ProjectConfig struct {
	Port        string
	DatabaseURL string
	DebugMode   bool
}

func LoadConfig() *ProjectConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	debugMode, err := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
	if err != nil {
		log.Fatalf("Invalid value for DEBUG_MODE: %v", err)
	}

	return &ProjectConfig{
		Port:        getEnv("PORT", "3001"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DebugMode:   debugMode,
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
