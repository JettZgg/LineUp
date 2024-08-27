package api

import (
	"time"
	"strconv"

	"github.com/JettZgg/LineUp/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/JettZgg/LineUp/internal/websocket"
)

func SetupRoutes(r *gin.Engine) {
	// Create a new hub
	hub := websocket.NewHub()
	go hub.Run()

	// Middleware to add hub to context
	r.Use(func(c *gin.Context) {
		c.Set("hub", hub)
		c.Next()
	})

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Update this with your client's URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.POST("/api/register", RegisterHandler)
	r.POST("/api/login", LoginHandler)

	// WebSocket route
	r.GET("/ws/:matchID", func(c *gin.Context) {
		matchID, _ := strconv.ParseInt(c.Param("matchID"), 10, 64)
		websocket.ServeWs(hub, c.Writer, c.Request, matchID)
	})

	// Authenticated routes
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// Match-related routes
		auth.POST("/create-match", CreateMatchHandler)
		auth.POST("/join-match/:matchID", JoinMatchHandler)
		auth.GET("/match/:matchID", GetMatchHandler)
		auth.GET("/match-history", GetMatchHistoryHandler)
		auth.GET("/match-replay/:matchID", GetMatchReplayHandler)

		// Add any other authenticated routes here
	}
}