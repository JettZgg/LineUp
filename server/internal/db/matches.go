// File: internal/db/matches.go
package db

import (
	"time"
)

type Match struct {
	ID         string    `json:"id"`
	Player1ID  int       `json:"player1_id"`
	Player2ID  int       `json:"player2_id"`
	BoardState string    `json:"board_state"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func CreateMatch(match *Match) error {
	_, err := DB.Exec(`
		INSERT INTO matches (id, player1_id, player2_id, board_state, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, match.ID, match.Player1ID, match.Player2ID, match.BoardState, match.Status, time.Now(), time.Now())
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

func GetMatchByID(id string) (*Match, error) {
	match := &Match{}
	err := DB.QueryRow("SELECT id, player1_id, player2_id, board_state, status, created_at, updated_at FROM matches WHERE id = $1", id).
		Scan(&match.ID, &match.Player1ID, &match.Player2ID, &match.BoardState, &match.Status, &match.CreatedAt, &match.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return match, nil
}
