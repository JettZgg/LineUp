// File: main.go
package main

import (
	"log"

	"github.com/JettZgg/LineUp/internal/config"
	"github.com/JettZgg/LineUp/internal/server"
	"github.com/JettZgg/LineUp/internal/utils"
	"github.com/JettZgg/LineUp/internal/utils/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize Snowflake nodes
	if err := utils.InitSnowflake(); err != nil {
		log.Fatalf("Failed to initialize Snowflake: %v", err)
	}

	// Create and start server
	srv := server.New(cfg, hub)
	log.Fatal(srv.Start())
}