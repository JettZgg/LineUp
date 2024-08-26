package main

import (
	"log"

	"github.com/JettZgg/LineUp/internal/config"
	"github.com/JettZgg/LineUp/internal/api"
	"github.com/JettZgg/LineUp/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Snowflake nodes
	if err := utils.InitSnowflake(); err != nil {
		log.Fatalf("Failed to initialize Snowflake: %v", err)
	}

	// Create and start API server
	apiServer := api.New(cfg)
	log.Fatal(apiServer.Start())
}