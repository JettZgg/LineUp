package websocket

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JettZgg/LineUp/internal/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
	hub    *Hub
	router *gin.Engine
}

func NewServer(cfg *config.Config, hub *Hub) *Server {
	router := gin.Default()
	server := &Server{
		config: cfg,
		hub:    hub,
		router: router,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.GET("/ws/:matchID", func(c *gin.Context) {
		matchID, err := strconv.ParseInt(c.Param("matchID"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
			return
		}
		ServeWs(s.hub, c.Writer, c.Request, matchID)
	})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.WebSocket.Port)
	fmt.Printf("WebSocket server starting on %s\n", addr)
	return s.router.Run(addr)
}