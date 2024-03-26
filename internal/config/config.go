package config

import (
	"errors"
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

func getenvInt(val string, def int) (int, error) {
	s := os.Getenv(val)
	if s == "" {
		return def, nil
	}
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.Join(errors.New("config .getenvInt(): "), err)
	}
	return result, nil
}

func NewConfig() (Config, error) {
	port, err := getenvInt("DB_PORT", 5432)
	return Config{
		DBPassword: getenv("DB_PASSWORD", "password"),
		DBName:     getenv("DB_NAME", "postgres"),
		DBHost:     getenv("DB_HOST", "localhost"),
		DBUser:     getenv("DB_USER", "postgres"),
		DBPort:     port,
	}, err
}
