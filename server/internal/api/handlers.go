// File: internal/api/handlers.go
package api

import (
	"log"
	"net/http"
	"time"

	"github.com/JettZgg/LineUp/internal/auth"
	"github.com/JettZgg/LineUp/internal/db"
	"github.com/JettZgg/LineUp/internal/game"
	"github.com/JettZgg/LineUp/internal/utils/websocket"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := auth.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func LoginHandler(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := auth.LoginUser(credentials.Username, credentials.Password)
	if err != nil {
		log.Printf("Login error for user %s: %v", credentials.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CreateMatchHandler(c *gin.Context) {
	var matchConfig game.MatchConfig
	if err := c.ShouldBindJSON(&matchConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	match, err := game.CreateMatch(matchConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create match"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"match":      match,
		"serverTime": time.Now().UTC(),
	})
}

func JoinMatchHandler(c *gin.Context) {
	matchID := c.Param("matchID")
	playerID := c.GetString("username") // Assuming the username is set in the context by the auth middleware

	if err := game.JoinMatch(matchID, playerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to join match"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined match successfully"})
}

func MakeMoveHandler(c *gin.Context) {
	matchID := c.Param("matchID")
	playerID := c.GetString("username")

	var moveRequest struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	if err := c.ShouldBindJSON(&moveRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get the hub from the context or wherever it's stored
	hub := c.MustGet("hub").(*websocket.Hub)

	result, err := game.MakeMove(hub, matchID, playerID, moveRequest.X, moveRequest.Y)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid move"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetMatchHandler(c *gin.Context) {
	matchID := c.Param("matchID")
	match, err := game.GetMatch(matchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"match":      match,
		"serverTime": time.Now().UTC(),
	})
}

func GetMatchHistoryHandler(c *gin.Context) {
	userID := c.GetString("username") // Assuming the username is set in the context by the auth middleware
	limit := 10                       // Or get this from query parameter

	matches, err := game.GetMatchHistory(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve match history"})
		return
	}

	c.JSON(http.StatusOK, matches)
}

func GetMatchReplayHandler(c *gin.Context) {
	matchID := c.Param("matchID")

	match, moves, err := game.GetMatchReplay(matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve match replay"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"match": match,
		"moves": moves,
	})
}
