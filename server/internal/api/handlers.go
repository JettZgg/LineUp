package api

import (
	"encoding/json"
	"net/http"

	"github.com/JettZgg/LineUp/internal/auth"
	"github.com/JettZgg/LineUp/internal/db"
	"github.com/JettZgg/LineUp/internal/game"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := auth.RegisterUser(&user); err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := auth.LoginUser(credentials.Username, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func CreateMatchHandler(w http.ResponseWriter, r *http.Request) {
	var matchConfig game.MatchConfig
	if err := json.NewDecoder(r.Body).Decode(&matchConfig); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	match, err := game.CreateMatch(matchConfig)
	if err != nil {
		http.Error(w, "Failed to create match", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(match)
}

func JoinMatchHandler(w http.ResponseWriter, r *http.Request) {
	var joinRequest struct {
		MatchID  string `json:"matchId"`
		PlayerID string `json:"playerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&joinRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := game.JoinMatch(joinRequest.MatchID, joinRequest.PlayerID); err != nil {
		http.Error(w, "Failed to join match", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Joined match successfully"})
}

func MakeMoveHandler(w http.ResponseWriter, r *http.Request) {
	var moveRequest struct {
		MatchID  string `json:"matchId"`
		PlayerID string `json:"playerId"`
		X        int    `json:"x"`
		Y        int    `json:"y"`
	}
	if err := json.NewDecoder(r.Body).Decode(&moveRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := game.MakeMove(moveRequest.MatchID, moveRequest.PlayerID, moveRequest.X, moveRequest.Y)
	if err != nil {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(result)
}
