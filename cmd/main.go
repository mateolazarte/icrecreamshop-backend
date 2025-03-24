package main

import (
	"icecreamshop/internal/api"
	"icecreamshop/internal/storage"
	"log"
	"os"
)

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = validateEnvironments()
	if err != nil {
		log.Fatal(err)
	}

	users := initialUsers()
	prices := initialPrices()
	flavors := initialFlavors()

	api_env := os.Getenv("API_ENV")
	var db storage.Storage
	if api_env == "development" || api_env == "production" {
		db = storage.NewDBStorage(flavors, users, prices)
	} else if api_env == "testing" {
		db = storage.NewMemoryStorage(flavors, users, prices)
	} else {
		log.Fatal("Invalid API_ENV")
	}

	log.Printf("Server running in %v mode\n", api_env)
	sv := api.NewServer(db)
	log.Fatal(sv.Start())
}
