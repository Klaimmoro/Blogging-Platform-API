package db

import "os"

type DBConfig struct {
	Host     string `json:"DB_HOST"`
	Port     string `json:"DB_PORT"`
	User     string `json:"DB_USER"`
	Password string `json:"DB_PASSWORD"`
	Name     string `json:"DB_NAME"`
}

func SetConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}
