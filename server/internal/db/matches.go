package db

import (
	"time"
)

type Match struct {
	ID          string    `json:"id"`
	Player1ID   string    `json:"player1_id"`
	Player2ID   string    `json:"player2_id"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time,omitempty"`
	Winner      string    `json:"winner,omitempty"`
	BoardWidth  int       `json:"board_width"`
	BoardHeight int       `json:"board_height"`
	WinLength   int       `json:"win_length"`
}

func CreateMatch(match *Match) error {
	_, err := DB.Exec(`
        INSERT INTO matches (id, player1_id, status, start_time, board_width, board_height, win_length)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, match.ID, match.Player1ID, match.Status, match.StartTime, match.BoardWidth, match.BoardHeight, match.WinLength)
	return err
}

func UpdateMatch(match *Match) error {
	_, err := DB.Exec(`
        UPDATE matches
        SET player2_id = $2, status = $3, end_time = $4, winner = $5
        WHERE id = $1
    `, match.ID, match.Player2ID, match.Status, match.EndTime, match.Winner)
	return err
}

func GetMatchByID(matchID string) (*Match, error) {
	match := &Match{}
	err := DB.QueryRow(`
        SELECT id, player1_id, player2_id, status, start_time, end_time, winner, board_width, board_height, win_length
        FROM matches WHERE id = $1
    `, matchID).Scan(
		&match.ID, &match.Player1ID, &match.Player2ID, &match.Status, &match.StartTime,
		&match.EndTime, &match.Winner, &match.BoardWidth, &match.BoardHeight, &match.WinLength,
	)
	if err != nil {
		return nil, err
	}
	return match, nil
}

func GetRecentMatchesByUser(userID string, limit int) ([]Match, error) {
	rows, err := DB.Query(`
        SELECT id, player1_id, player2_id, status, start_time, end_time, winner, board_width, board_height, win_length 
        FROM matches 
        WHERE player1_id = $1 OR player2_id = $1 
        ORDER BY start_time DESC 
        LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.ID, &match.Player1ID, &match.Player2ID, &match.Status, &match.StartTime,
			&match.EndTime, &match.Winner, &match.BoardWidth, &match.BoardHeight, &match.WinLength); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, nil
}
