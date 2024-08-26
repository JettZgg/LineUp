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
	"github.com/JettZgg/LineUp/internal/websocket"
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

	c.JSON(http.StatusOK, gin.H{"token": token, "userID": user.UID, "username": user.Username})
}

func CreateMatchHandler(c *gin.Context) {
	uid := c.GetInt64("uid")

	match, err := game.CreateMatch(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	player, err := db.GetUserByID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get player info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"match": match,
		"player": gin.H{
			"id":       strconv.FormatInt(player.UID, 10),
			"username": player.Username,
		},
	})

	// Create a new room in the WebSocket hub
	hub := c.MustGet("hub").(*websocket.Hub)
	hub.Rooms[match.MID] = &websocket.Room{
		ID:      match.MID,
		Clients: make(map[*websocket.Client]bool),
	}
}

func JoinMatchHandler(c *gin.Context) {
	matchID, err := strconv.ParseInt(c.Param("matchID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}
	uid := c.GetInt64("uid")

	hub := c.MustGet("hub").(*websocket.Hub)
	broadcastFunc := func(matchID int64, message []byte) {
		hub.BroadcastToMatch(matchID, message)
	}

	if err := game.JoinMatch(broadcastFunc, matchID, uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined match successfully"})
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

	replay, err := game.GetMatchReplay(matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve match replay"})
		return
	}

	c.JSON(http.StatusOK, replay)
}
