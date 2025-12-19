package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBURL         string
	DBMaxConns    int32
	DBMinConns    int32
	DBMaxWaitTime time.Duration
}

func NewConfig() *Config {
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("No config.env file found, reading from environment")
	}

	return &Config{
		Port:          os.Getenv("PORT"),
		DBURL:         os.Getenv("DB_URL"),
		DBMaxConns:    int32(getEnvInt("DB_MAX_CONNS", 100)),
		DBMinConns:    int32(getEnvInt("DB_MIN_CONNS", 10)),
		DBMaxWaitTime: getEnvDuration("DB_MAX_CONN_WAIT", 30*time.Second),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if d, err := time.ParseDuration(valueStr); err == nil {
		return d
	}
	return fallback
}
