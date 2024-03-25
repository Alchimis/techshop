package config

import "os"

type Config struct {
	DBPassword, DBName, DBHost string
}

func Getenv(val, def string) string {
	s := os.Getenv(val)
	if s == "" {
		return def
	}
	return s
}

func NewConfig() Config {
	return Config{
		DBPassword: Getenv("DB_PASSWORD", "password"),
		DBName:     Getenv("DB_NAME", "postgres"),
		DBHost:     Getenv("DB_HOST", "localhost"),
	}
}
