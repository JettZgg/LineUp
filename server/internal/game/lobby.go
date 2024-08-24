package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/JettZgg/LineUp/internal/db"
	"github.com/JettZgg/LineUp/internal/utils"
)

func CreateMatch(playerID int64, config MatchConfig) (*Match, error) {
	matchID := utils.GenerateMID()
	match := &Match{
		MID:       matchID,
		Player1ID: playerID,
		Status:    "waiting",
		Config:    config,
		StartTime: time.Now().UTC(),
	}

	dbMatch := &db.Match{
		MID:       match.MID,
		Player1ID: playerID,
		Status:    match.Status,
		StartTime: match.StartTime,
	}

	if err := db.CreateMatch(dbMatch); err != nil {
		return nil, fmt.Errorf("failed to create match in database: %w", err)
	}

	matches[match.MID] = match
	return match, nil
}

func JoinMatch(broadcastFunc func(int64, []byte), matchID int64, playerID int64) error {
	match, exists := matches[matchID]
	if !exists {
		return errors.New("match not found")
	}

	// Always update the player, whether they're joining or rejoining
	if match.Player1ID == 0 {
		match.Player1ID = playerID
	} else if match.Player2ID == 0 && match.Player1ID != playerID {
		match.Player2ID = playerID
	}

	match.Status = "waiting_ready"

	// Update the match in the database
	dbMatch := &db.Match{
		MID:       match.MID,
		Player1ID: match.Player1ID,
		Player2ID: match.Player2ID,
		Status:    match.Status,
	}
	if err := db.UpdateMatch(dbMatch); err != nil {
		log.Printf("Error updating match in database: %v", err)
		return fmt.Errorf("failed to update match in database: %w", err)
	}

	log.Printf("Player %d joined/rejoined match %d", playerID, matchID)

	// Broadcast updated game info to all players
	return broadcastGameInfo(broadcastFunc, matchID)
}

func GetMatchHistory(userID int64, limit int) ([]db.Match, error) {
	return db.GetRecentMatchesByUser(userID, limit)
}

func GetMatchReplay(matchID int64) (*db.Match, error) {
	return db.GetMatchByID(matchID)
}

func LeaveMatch(broadcastFunc func(int64, []byte), matchID int64, playerID int64) error {
	match, exists := matches[matchID]
	if !exists {
		return errors.New("match not found")
	}

	if match.Player1ID == playerID {
		if match.Player2ID != 0 {
			match.Player1ID = match.Player2ID
			match.Player2ID = 0
			match.Player1Ready = match.Player2Ready
			match.Player2Ready = false
		} else {
			delete(matches, matchID)
			return db.DeleteMatch(matchID)
		}
	} else if match.Player2ID == playerID {
		match.Player2ID = 0
		match.Player2Ready = false
	} else {
		return errors.New("player not in this match")
	}

	match.Status = "waiting"

	dbMatch := &db.Match{
		MID:       match.MID,
		Player1ID: match.Player1ID,
		Player2ID: match.Player2ID,
		Status:    match.Status,
	}
	if err := db.UpdateMatch(dbMatch); err != nil {
		log.Printf("Error updating match in database: %v", err)
		return fmt.Errorf("failed to update match in database: %w", err)
	}

	gameInfo, err := GetGameInfo(matchID)
	if err != nil {
		return fmt.Errorf("failed to get game info: %w", err)
	}

	msgBytes, _ := json.Marshal(gameInfo)
	broadcastFunc(matchID, msgBytes)

	return nil
}
