package main

import (
	"errors"

	"os"

	"strings"

	"github.com/joho/godotenv"
)

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}

// validateEnvironments checks if all env variables were inputted.
func validateEnvironments() error {
	// API
	if strings.TrimSpace(os.Getenv("API_ENV")) == "" {
		return errors.New("SERVER_PORT env is needed")
	}
	if strings.TrimSpace(os.Getenv("API_PORT")) == "" {
		return errors.New(" SERVER_PORT env is needed")
	}

	// jwt
	if strings.TrimSpace(os.Getenv("JWT_SECRET")) == "" {
		return errors.New("JWT_SECRET env is needed")
	}

	// Database
	if strings.TrimSpace(os.Getenv("DB_USER")) == "" {
		return errors.New("DB_USER env is needed")
	}
	if strings.TrimSpace(os.Getenv("DB_PASSWORD")) == "" {
		return errors.New("DB_PASSWORD env is needed")
	}
	if strings.TrimSpace(os.Getenv("DB_HOST")) == "" {
		return errors.New("DB_HOST env is needed")
	}
	if strings.TrimSpace(os.Getenv("DB_PORT")) == "" {
		return errors.New("DB_PORT env is needed")
	}
	if strings.TrimSpace(os.Getenv("DB_NAME")) == "" {
		return errors.New("DB_NAME env is needed")
	}

	// Test Database
	if strings.TrimSpace(os.Getenv("TEST_DB_USER")) == "" {
		return errors.New("TEST_DB_USER env is needed")
	}
	if strings.TrimSpace(os.Getenv("TEST_DB_PASSWORD")) == "" {
		return errors.New("TEST_DB_PASSWORD env is needed")
	}
	if strings.TrimSpace(os.Getenv("TEST_DB_HOST")) == "" {
		return errors.New("TEST_DB_HOST env is needed")
	}
	if strings.TrimSpace(os.Getenv("TEST_DB_PORT")) == "" {
		return errors.New("TEST_DB_PORT env is needed")
	}
	if strings.TrimSpace(os.Getenv("TEST_DB_NAME")) == "" {
		return errors.New("TEST_DB_NAME env is needed")
	}

	//external services
	if strings.TrimSpace(os.Getenv("MP_ACCESS_TOKEN")) == "" {
		return errors.New("MP_ACCESS_TOKEN env is needed")
	}

	return nil
}
