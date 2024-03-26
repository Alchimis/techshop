package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBPassword, DBName, DBHost, DBUser string
	DBPort                             int
}

func getenv(val, def string) string {
	s := os.Getenv(val)
	if s == "" {
		return def
	}
	return s
}

func getenvInt(val string, def int) int {
	s := os.Getenv(val)
	if s == "" {
		return def
	}
	result, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return result
}

func NewConfig() Config {
	return Config{
		DBPassword: getenv("DB_PASSWORD", "password"),
		DBName:     getenv("DB_NAME", "postgres"),
		DBHost:     getenv("DB_HOST", "localhost"),
		DBUser:     getenv("DB_USER", "postgres"),
		DBPort:     getenvInt("DB_PORT", 5432),
	}
}
