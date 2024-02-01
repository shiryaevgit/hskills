package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Environment struct {
	HTTPPort         string
	DatabaseFilepath string
}

func LoadEnv() (*Environment, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loadEnv: %v", err)
	}
	env := &Environment{
		HTTPPort:         os.Getenv("HTTP_PORT"),
		DatabaseFilepath: os.Getenv("DATABASE_FILEPATH"),
	}
	return env, nil
}
