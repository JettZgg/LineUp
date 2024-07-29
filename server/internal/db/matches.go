// File: internal/db/matches.go
package db

import (
	"database/sql"
	"time"
)

type Match struct {
	ID         string       `json:"id"`
	Player1ID  int          `json:"player1_id"`
	Player2ID  int          `json:"player2_id"`
	BoardState string       `json:"board_state"`
	Status     string       `json:"status"`
	StartTime  sql.NullTime `json:"start_time"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

func CreateMatch(match *Match) error {
	now := time.Now()
	_, err := DB.Exec(`
		INSERT INTO matches (id, player1_id, player2_id, board_state, status, start_time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, match.ID, match.Player1ID, match.Player2ID, match.BoardState, match.Status, sql.NullTime{}, now, now)
	return err
}

func StartMatch(matchID string) error {
	now := time.Now()
	_, err := DB.Exec(`
		UPDATE matches
		SET start_time = $1, status = 'ongoing', updated_at = $2
		WHERE id = $3 AND start_time IS NULL
	`, now, now, matchID)
	return err
}

func UpdateMatch(match *Match) error {
	_, err := DB.Exec(`
		UPDATE matches
		SET player2_id = $2, board_state = $3, status = $4, updated_at = $5
		WHERE id = $1
	`, match.ID, match.Player2ID, match.BoardState, match.Status, time.Now())
	return err
}
