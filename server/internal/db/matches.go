package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Move struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Moves []Move

func (m Moves) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Moves) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &m)
}

type Match struct {
    MID               int64  `json:"id"`
    Player1ID         int64  `json:"player1_id"`
    Player2ID         int64  `json:"player2_id"`
    Winner            int64  `json:"winner,omitempty"`
    FirstMovePlayerID int64  `json:"first_move_player_id"`
    Moves             string `json:"moves"`
    Date              string `json:"date"`
}

func CreateMatch(match *Match) error {
    _, err := DB.Exec(`
        INSERT INTO matches (id, player1_id, first_move_player_id, date)
        VALUES ($1, $2, $3, $4)
    `, match.MID, match.Player1ID, match.FirstMovePlayerID, match.Date)
    return err
}

func UpdateMatch(match *Match) error {
    _, err := DB.Exec(`
        UPDATE matches
        SET player2_id = $2, winner = $3, moves = $4
        WHERE id = $1
    `, match.MID, match.Player2ID, match.Winner, match.Moves)
    return err
}

func GetMatchByID(matchID int64) (*Match, error) {
    if DB == nil {
        return nil, errors.New("database connection is not initialized")
    }
    match := &Match{}
    err := DB.QueryRow(`
        SELECT id, player1_id, player2_id, winner, first_move_player_id, moves, date
        FROM matches WHERE id = $1
    `, matchID).Scan(
        &match.MID, &match.Player1ID, &match.Player2ID, &match.Winner,
        &match.FirstMovePlayerID, &match.Moves, &match.Date,
    )
    if err != nil {
        return nil, err
    }
    return match, nil
}

func GetRecentMatchesByUser(userID int64, limit int) ([]Match, error) {
    rows, err := DB.Query(`
        SELECT id, player1_id, player2_id, winner, first_move_player_id, moves, date
        FROM matches 
        WHERE player1_id = $1 OR player2_id = $1 
        ORDER BY date DESC 
        LIMIT $2
    `, userID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var matches []Match
    for rows.Next() {
        var match Match
        if err := rows.Scan(&match.MID, &match.Player1ID, &match.Player2ID, &match.Winner,
            &match.FirstMovePlayerID, &match.Moves, &match.Date); err != nil {
            return nil, err
        }
        matches = append(matches, match)
    }
    return matches, nil
}

func DeleteMatch(matchID int64) error {
    _, err := DB.Exec("DELETE FROM matches WHERE id = $1", matchID)
    if err != nil {
        return fmt.Errorf("failed to delete match: %w", err)
    }
    return nil
}