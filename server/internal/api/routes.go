package api

import (
	"time"

	"github.com/JettZgg/LineUp/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
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

	// Authenticated routes
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// Match-related routes
		auth.POST("/create-match", CreateMatchHandler)
		auth.POST("/join-match/:matchID", JoinMatchHandler)
		auth.POST("/make-move/:matchID", MakeMoveHandler)
		auth.GET("/match/:matchID", GetMatchHandler)
		auth.GET("/match-history", GetMatchHistoryHandler)
		auth.GET("/match-replay/:matchID", GetMatchReplayHandler)

		// Add any other authenticated routes here
	}
}