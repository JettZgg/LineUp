package match

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/JettZgg/LineUp/internal/db"
)

func GetMatchInfo(matchID int64) (map[string]interface{}, error) {
	match, err := GetMatch(matchID)
	if err != nil {
		log.Printf("Error getting match %d: %v", matchID, err)
		return nil, err
	}

	log.Printf("Match data: %+v", match)

	players := []map[string]interface{}{}

	if match.Player1ID != 0 {
		player1, err := db.GetUserByID(match.Player1ID)
		if err != nil {
			log.Printf("Error getting player1 (ID: %d): %v", match.Player1ID, err)
			return nil, err
		}
		players = append(players, map[string]interface{}{
			"id":       player1.UID,
			"username": player1.Username,
		})
	}

	if match.Player2ID != 0 {
		player2, err := db.GetUserByID(match.Player2ID)
		if err != nil {
			log.Printf("Error getting player2 (ID: %d): %v", match.Player2ID, err)
			return nil, err
		}
		players = append(players, map[string]interface{}{
			"id":       player2.UID,
			"username": player2.Username,
		})
	}

	return map[string]interface{}{
		"type":    "matchInfo",
		"matchId": matchID,
		"players": players,
		"moves":   match.Moves,
		"winner":  match.Winner,
		"date":    match.Date,
	}, nil
}

func broadcastMatchInfo(broadcastFunc func(int64, []byte), matchID int64) error {
	matchInfo, err := GetMatchInfo(matchID)
	if err != nil {
		return fmt.Errorf("failed to get match info: %w", err)
	}

	msgBytes, err := json.Marshal(matchInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal match info: %w", err)
	}

	broadcastFunc(matchID, msgBytes)

	return nil
}

func updateMatchInDatabase(match *Match) error {
	dbMatch := &db.Match{
		MID:       match.MID,
		Player2ID: match.Player2ID,
		Winner:    match.Winner,
		Moves:     match.Moves,
	}
	return db.UpdateMatch(dbMatch)
}
