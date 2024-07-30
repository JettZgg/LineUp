package db

import (
	"log"
	"time"
)

type Match struct {
	MID         int64     `json:"id"`
	Player1ID   int64     `json:"player1_id"`
	Player2ID   int64     `json:"player2_id"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time,omitempty"`
	Winner      int64     `json:"winner,omitempty"`
	BoardWidth  int       `json:"board_width"`
	BoardHeight int       `json:"board_height"`
	WinLength   int       `json:"win_length"`
}

func CreateMatch(match *Match) error {
	_, err := DB.Exec(`
        INSERT INTO matches (id, player1_id, status, start_time, board_width, board_height, win_length)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, match.MID, match.Player1ID, match.Status, match.StartTime, match.BoardWidth, match.BoardHeight, match.WinLength)
	return err
}

func UpdateMatch(match *Match) error {
	_, err := DB.Exec(`
        UPDATE matches
        SET player2_id = $2, status = $3, end_time = $4, winner = $5
        WHERE id = $1
    `, match.MID, match.Player2ID, match.Status, match.EndTime, match.Winner)
	return err
}

func GetMatchByID(matchID int64) (*Match, error) {
	match := &Match{}
	err := DB.QueryRow(`
        SELECT id, player1_id, player2_id, status, start_time, end_time, winner, board_width, board_height, win_length
        FROM matches WHERE id = $1
    `, matchID).Scan(
		&match.MID, &match.Player1ID, &match.Player2ID, &match.Status, &match.StartTime,
		&match.EndTime, &match.Winner, &match.BoardWidth, &match.BoardHeight, &match.WinLength,
	)
	if err != nil {
		return nil, err
	}
	return match, nil
}

func GetRecentMatchesByUser(userID int64, limit int) ([]Match, error) {
	rows, err := DB.Query(`
        SELECT id, player1_id, player2_id, status, start_time, end_time, winner, board_width, board_height, win_length 
        FROM matches 
        WHERE player1_id = $1 OR player2_id = $1 
        ORDER BY start_time DESC 
        LIMIT $2
    `, userID, limit)
	if err != nil {
		log.Printf("Error querying recent matches: %v", err)
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.MID, &match.Player1ID, &match.Player2ID, &match.Status, &match.StartTime,
			&match.EndTime, &match.Winner, &match.BoardWidth, &match.BoardHeight, &match.WinLength); err != nil {
			log.Printf("Error scanning match row: %v", err)
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, nil
}
