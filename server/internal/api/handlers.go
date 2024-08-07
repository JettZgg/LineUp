// File: internal/api/handlers.go
package api

import (
	"log"
	"net/http"
	"strconv"
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

	user, token, err := auth.LoginUser(credentials.Username, credentials.Password)
	if err != nil {
		log.Printf("Login error for user %s: %v", credentials.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Set the user UID in the context
	c.Set("userID", user.UID)

	c.JSON(http.StatusOK, gin.H{"token": token, "userID": user.UID})
}

func CreateMatchHandler(c *gin.Context) {
	var matchConfig game.MatchConfig
	if err := c.ShouldBindJSON(&matchConfig); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get the user ID from the authenticated context
	uid := c.GetInt64("uid") // Make sure your auth middleware sets this

	log.Printf("Attempting to create match with config: %+v for user: %d", matchConfig, uid)

	match, err := game.CreateMatch(matchConfig, uid)
	if err != nil {
		log.Printf("Error creating match: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create match", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"match":      match,
		"serverTime": time.Now().UTC(),
	})
}

func JoinMatchHandler(c *gin.Context) {
	matchID, err := strconv.ParseInt(c.Param("matchID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}
	uid := c.GetInt64("uid")

	if err := game.JoinMatch(matchID, uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to join match"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined match successfully"})
}

func MakeMoveHandler(c *gin.Context) {
	matchID, err := strconv.ParseInt(c.Param("matchID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}
	uid := c.GetInt64("uid")

	var moveRequest struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	if err := game.JoinMatch(matchID, uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get the hub from the context or wherever it's stored
	hub := c.MustGet("hub").(*websocket.Hub)

	result, err := game.MakeMove(hub, matchID, uid, moveRequest.X, moveRequest.Y)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid move"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetMatchHandler(c *gin.Context) {
	matchID, err := strconv.ParseInt(c.Param("matchID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

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
	uid := c.GetInt64("uid") // Assuming the username is set in the context by the auth middleware
	limit := 10              // Or get this from query parameter

	matches, err := game.GetMatchHistory(uid, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve match history"})
		return
	}

	c.JSON(http.StatusOK, matches)
}

func GetMatchReplayHandler(c *gin.Context) {
	matchID, err := strconv.ParseInt(c.Param("matchID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

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
