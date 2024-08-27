package match

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/JettZgg/LineUp/internal/db"
	"github.com/JettZgg/LineUp/internal/utils"
)

func CreateMatch(playerID int64) (*Match, error) {
	matchID := utils.GenerateMID()
	currentDate := time.Now().Format("2006-01-02")
	match := &Match{
		MID:               matchID,
		Player1ID:         playerID,
		FirstMovePlayerID: playerID,
		Date:              currentDate,
	}

	dbMatch := &db.Match{
		MID:               match.MID,
		Player1ID:         playerID,
		FirstMovePlayerID: playerID,
		Date:              currentDate,
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
		// If the match is not in memory, try to fetch it from the database
		dbMatch, err := db.GetMatchByID(matchID)
		if err != nil {
			return fmt.Errorf("match not found: %w", err)
		}
		match = &Match{
			MID:               dbMatch.MID,
			Player1ID:         dbMatch.Player1ID,
			Player2ID:         dbMatch.Player2ID,
			FirstMovePlayerID: dbMatch.FirstMovePlayerID,
			Winner:            dbMatch.Winner,
			Moves:             dbMatch.Moves,
			Date:              dbMatch.Date,
		}
		matches[matchID] = match
	}

	if match.Player1ID == 0 {
		match.Player1ID = playerID
	} else if match.Player2ID == 0 && match.Player1ID != playerID {
		match.Player2ID = playerID
	} else {
		return errors.New("match is full or player already joined")
	}

	dbMatch := &db.Match{
		MID:       match.MID,
		Player1ID: match.Player1ID,
		Player2ID: match.Player2ID,
	}
	if err := db.UpdateMatch(dbMatch); err != nil {
		log.Printf("Error updating match in database: %v", err)
		return fmt.Errorf("failed to update match in database: %w", err)
	}

	log.Printf("Player %d joined match %d", playerID, matchID)

	return broadcastMatchInfo(broadcastFunc, matchID)
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
		match.Player1ID = match.Player2ID
		match.Player2ID = 0
	} else if match.Player2ID == playerID {
		match.Player2ID = 0
	} else {
		return errors.New("player not in this match")
	}

	dbMatch := &db.Match{
		MID:       match.MID,
		Player1ID: match.Player1ID,
		Player2ID: match.Player2ID,
	}
	if err := db.UpdateMatch(dbMatch); err != nil {
		log.Printf("Error updating match in database: %v", err)
		return fmt.Errorf("failed to update match in database: %w", err)
	}

	matchInfo, err := GetMatchInfo(matchID)
	if err != nil {
		return fmt.Errorf("failed to get match info: %w", err)
	}

	msgBytes, _ := json.Marshal(matchInfo)
	broadcastFunc(matchID, msgBytes)

	return nil
}
