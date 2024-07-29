package server

import (
	"fmt"
	"net/http"

	"github.com/JettZgg/LineUp/internal/api"
	"github.com/JettZgg/LineUp/internal/config"
	"github.com/JettZgg/LineUp/internal/db"
	"github.com/JettZgg/LineUp/internal/websocket"
)

type Server struct {
	config *config.Config
	router http.Handler
}

func New(cfg *config.Config, hub *websocket.Hub) (*Server, error) {
	if err := db.Initialize(cfg.Database); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	router := api.SetupRoutes(hub)

	return &Server{
		config: cfg,
		router: router,
	}, nil
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}