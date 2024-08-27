package main

import (
	"log"

	"github.com/DimTur/chat-websocket-go/internal/pkg/client/app"
	"github.com/DimTur/chat-websocket-go/internal/pkg/client/internal/config"
)

func main() {
	_, err := config.LoadConfig("./internal/config/config.toml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app.Run()
}
