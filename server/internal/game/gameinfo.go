package game

import (
	"log"

	"github.com/JettZgg/LineUp/internal/db"
)

func GetGameInfo(matchID int64) (map[string]interface{}, error) {
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
			"id":     player1.UID,
			"username": player1.Username,
			"ready":   match.Player1Ready,
		})
	}

	if match.Player2ID != 0 {
		player2, err := db.GetUserByID(match.Player2ID)
		if err != nil {
			log.Printf("Error getting player2 (ID: %d): %v", match.Player2ID, err)
			return nil, err
		}
		players = append(players, map[string]interface{}{
			"id":     player2.UID,
			"username": player2.Username,
			"ready":   match.Player2Ready,
		})
	}

	log.Printf("GetGameInfo for match %d: players=%v, config=%v", matchID, players, match.Config)

	return map[string]interface{}{
		"type":    "gameInfo",
		"matchId": matchID,
		"players": players,
		"config":  match.Config,
	}, nil
}