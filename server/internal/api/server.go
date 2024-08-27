package api

import (
	"fmt"

	"github.com/JettZgg/LineUp/internal/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
	router *gin.Engine
}

func New(cfg *config.Config) *Server {
	router := gin.Default()
	SetupRoutes(router)

	return &Server{
		config: cfg,
		router: router,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	fmt.Printf("API server starting on %s\n", addr)
	return s.router.Run(addr)
}
