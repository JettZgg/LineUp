package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/JettZgg/LineUp/internal/db"
	"golang.org/x/exp/rand"
)

func UpdateGameConfig(broadcastFunc func(int64, []byte), matchID int64, playerID int64, config MatchConfig) error {
	match, exists := matches[matchID]
	if !exists {
		return errors.New("match not found")
	}
	if match.Player1ID != playerID {
		return errors.New("only the match owner can update the configuration")
	}

	match.Config = config

	// Update the match in the database
	dbMatch := &db.Match{
		MID:         match.MID,
		BoardWidth:  config.BoardWidth,
		BoardHeight: config.BoardHeight,
		WinLength:   config.WinLength,
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

func StartMatch(broadcastFunc func(int64, []byte), matchID int64, playerID int64) error {
	match, exists := matches[matchID]
	if !exists {
		return errors.New("match not found")
	}
	if match.Player1ID != playerID {
		return errors.New("only the match owner can start the match")
	}
	if match.Player2ID == 0 {
		return errors.New("waiting for second player to join")
	}
	if !match.Player1Ready || !match.Player2Ready {
		return errors.New("all players must be ready to start the match")
	}

	match.Status = "ongoing"
	match.StartTime = time.Now()
	match.FirstMovePlayerID = determineFirstPlayer(match.Player1ID, match.Player2ID)

	// Update the match in the database
	dbMatch := &db.Match{
		MID:               match.MID,
		Status:            match.Status,
		StartTime:         match.StartTime,
		FirstMovePlayerID: match.FirstMovePlayerID,
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

func SetPlayerReady(broadcastFunc func(int64, []byte), matchID int64, playerID int64, isReady bool) error {
	match, exists := matches[matchID]
	if !exists {
		return errors.New("match not found")
	}

	if match.Player1ID == playerID {
		match.Player1Ready = isReady
	} else if match.Player2ID == playerID {
		match.Player2Ready = isReady
	} else {
		return errors.New("player not in this match")
	}

	return broadcastGameInfo(broadcastFunc, matchID)
}

func determineFirstPlayer(player1ID, player2ID int64) int64 {
	if rand.Intn(2) == 0 {
		return player1ID
	}
	return player2ID
}
