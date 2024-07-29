// File: internal/server/server.go
package server

import (
	"fmt"

	"github.com/JettZgg/LineUp/internal/api"
	"github.com/JettZgg/LineUp/internal/config"
	"github.com/JettZgg/LineUp/internal/db"
	"github.com/JettZgg/LineUp/internal/utils/websocket"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
	router *gin.Engine
}

func New(cfg *config.Config, hub *websocket.Hub) *Server {
	if err := db.Initialize(cfg.Database); err != nil {
		panic(fmt.Errorf("failed to initialize database: %w", err))
	}

	router := gin.Default()
	api.SetupRoutes(router, hub)

	return &Server{
		config: cfg,
		router: router,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	return s.router.Run(addr)
}
