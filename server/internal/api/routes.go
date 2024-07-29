// File: internal/api/routes.go
package api

import (
	"github.com/JettZgg/LineUp/internal/middleware"
	"github.com/JettZgg/LineUp/internal/utils/websocket"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, hub *websocket.Hub) {
	// Public routes
	r.POST("/api/register", RegisterHandler)
	r.POST("/api/login", LoginHandler)

	// Authenticated routes
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// Match-related routes
		auth.POST("/create-match", CreateMatchHandler)
		auth.POST("/join-match/:matchID", JoinMatchHandler)
		auth.POST("/make-move/:matchID", MakeMoveHandler)
		auth.GET("/match/:matchID", GetMatchHandler) // New route for getting match details

		// Add any other authenticated routes here
	}

	// WebSocket route
	r.GET("/ws/:matchID", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})
}
