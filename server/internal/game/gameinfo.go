package game

import (
	"github.com/JettZgg/LineUp/internal/db"
)

func GetGameInfo(matchID int64) (map[string]interface{}, error) {
	match, err := GetMatch(matchID)
	if err != nil {
		return nil, err
	}

	player1, err := db.GetUserByID(match.Player1ID)
	if err != nil {
		return nil, err
	}

	var player2 *db.User
	if match.Player2ID != 0 {
		player2, err = db.GetUserByID(match.Player2ID)
		if err != nil {
			return nil, err
		}
	}

	players := []map[string]interface{}{
		{"id": player1.UID, "username": player1.Username},
	}
	if player2 != nil {
		players = append(players, map[string]interface{}{"id": player2.UID, "username": player2.Username})
	}

	return map[string]interface{}{
		"type":    "gameInfo",
		"matchId": matchID,
		"players": players,
		"config":  match.Config,
	}, nil
}
