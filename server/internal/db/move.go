package db

import (
	"time"
)

type Move struct {
	ID         int       `json:"id"`
	MatchID    int64     `json:"match_id"`
	PlayerID   int64     `json:"player_id"`
	X          int       `json:"x"`
	Y          int       `json:"y"`
	MoveNumber int       `json:"move_number"`
	Timestamp  time.Time `json:"timestamp"`
}

func CreateMove(move *Move) error {
	_, err := DB.Exec(`
        INSERT INTO moves (match_id, player_id, x, y, move_number, timestamp)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, move.MatchID, move.PlayerID, move.X, move.Y, move.MoveNumber, move.Timestamp)
	return err
}

func GetMovesByMatchID(matchID int64) ([]Move, error) {
	rows, err := DB.Query("SELECT id, match_id, player_id, x, y, move_number, timestamp FROM moves WHERE match_id = $1 ORDER BY move_number", matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var moves []Move
	for rows.Next() {
		var move Move
		if err := rows.Scan(&move.ID, &move.MatchID, &move.PlayerID, &move.X, &move.Y, &move.MoveNumber, &move.Timestamp); err != nil {
			return nil, err
		}
		moves = append(moves, move)
	}
	return moves, nil
}
