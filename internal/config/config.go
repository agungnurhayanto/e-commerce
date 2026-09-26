package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName        string
	AppPort        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	JWTSecret      string
	JWTExpireHours int
}

func Load() (*Config, error) {
	err :=
		godotenv.Load()
	if err != nil {
		return nil, err
	}

	jwtExpireHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))

	if err != nil {
		return nil, err
	}

	return &Config{

		AppName:        os.Getenv("APP_NAME"),
		AppPort:        os.Getenv("APP_PORT"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		DBSSLMode:      os.Getenv("DB_SSLMODE"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpireHours: jwtExpireHours,
	}, nil
}
