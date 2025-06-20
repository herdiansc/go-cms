package config

import (
	"github.com/joho/godotenv"
)

func LoadEnv(file string) {
	err := godotenv.Load(file)
	if err != nil {
		panic("Error loading .env file")
	}
}
