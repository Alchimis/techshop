package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBPassword, DBName, DBHost, DBUser string
}

func Getenv(val, def string) string {
	s := os.Getenv(val)
	if s == "" {
		return def
	}
	return s
}

func GetenvInt(val string, def int) int {
	s := os.Getenv(val)
	if s == "" {
		return def
	}
	return strconv.Atoi()
}

func NewConfig() Config {
	return Config{
		DBPassword: Getenv("DB_PASSWORD", "password"),
		DBName:     Getenv("DB_NAME", "postgres"),
		DBHost:     Getenv("DB_HOST", "localhost"),
		DBUser:     Getenv("DB_USER", "postgres"),
	}
}
