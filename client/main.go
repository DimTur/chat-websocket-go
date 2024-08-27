package main

import (
	"client/app"
	"client/internal/config"
	"log"
)

func main() {
	_, err := config.LoadConfig("./internal/config/config.toml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app.Run()
}
