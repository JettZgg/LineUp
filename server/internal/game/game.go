package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/JettZgg/LineUp/internal/db"
)

type MatchConfig struct {
	BoardWidth  int `json:"boardWidth"`
	BoardHeight int `json:"boardHeight"`
	WinLength   int `json:"winLength"`
}

type Match struct {
	MID               int64       `json:"id"`
	Board             [][]string  `json:"board"`
	Player1ID         int64       `json:"player1Id"`
	Player2ID         int64       `json:"player2Id"`
	Status            string      `json:"status"`
	Config            MatchConfig `json:"config"`
	StartTime         time.Time   `json:"startTime"`
	EndTime           time.Time   `json:"endTime"`
	Winner            int64       `json:"winner"`
	FirstMovePlayerID int64       `json:"firstMovePlayerId"`
	Moves             []db.Move   `json:"moves"`
	Player1Ready      bool        `json:"player1Ready"`
	Player2Ready      bool        `json:"player2Ready"`
}

func (m Match) MarshalJSON() ([]byte, error) {
	type Alias Match
	return json.Marshal(&struct {
		MID               string `json:"id"`
		Player1ID         string `json:"player1Id"`
		Player2ID         string `json:"player2Id"`
		FirstMovePlayerID string `json:"firstMovePlayerId"`
		Alias
	}{
		MID:               strconv.FormatInt(m.MID, 10),
		Player1ID:         strconv.FormatInt(m.Player1ID, 10),
		Player2ID:         strconv.FormatInt(m.Player2ID, 10),
		FirstMovePlayerID: strconv.FormatInt(m.FirstMovePlayerID, 10),
		Alias:             Alias(m),
	})
}

var matches = make(map[int64]*Match)

func GetMatch(matchID int64) (*Match, error) {
	match, exists := matches[matchID]
	if !exists {
		// If not in memory, try to fetch from database
		dbMatch, err := db.GetMatchByID(matchID)
		if err != nil {
			return nil, fmt.Errorf("match not found: %w", err)
		}
		// Convert db.Match to game.Match
		match = &Match{
			MID:               dbMatch.MID,
			Player1ID:         dbMatch.Player1ID,
			Player2ID:         dbMatch.Player2ID,
			Status:            dbMatch.Status,
			StartTime:         dbMatch.StartTime,
			FirstMovePlayerID: dbMatch.FirstMovePlayerID,
			Config: MatchConfig{
				BoardWidth:  dbMatch.BoardWidth,
				BoardHeight: dbMatch.BoardHeight,
				WinLength:   dbMatch.WinLength,
			},
			Board: makeBoard(dbMatch.BoardWidth, dbMatch.BoardHeight),
			Moves: dbMatch.Moves,
		}
		matches[matchID] = match
	}
	return match, nil
}

func MakeMove(broadcastFunc func(int64, []byte), matchID int64, playerID int64, x, y int) ([]byte, error) {
    match, err := GetMatch(matchID)
    if err != nil {
        return nil, err
    }
    if match.Status != "ongoing" {
        return nil, errors.New("match is not ongoing")
    }
    if (match.Player1ID != playerID && match.Player2ID != playerID) || (match.Board[y][x] != "") {
        return nil, errors.New("invalid move")
    }

    symbol := "X"
    if playerID == match.Player2ID {
        symbol = "O"
    }
    match.Board[y][x] = symbol

    // Add the move to the match's Moves slice
    move := db.Move{X: x, Y: y}
    match.Moves = append(match.Moves, move)

    result := checkGameResult(match, x, y)

    // Create a message to broadcast
    updateMsg := map[string]interface{}{
        "type":    "moveUpdate",
        "matchID": matchID,
        "board":   match.Board,
        "result":  result,
        "lastMove": map[string]interface{}{
            "playerID": playerID,
            "x":        x,
            "y":        y,
        },
    }

    msgBytes, err := json.Marshal(updateMsg)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal move update: %w", err)
    }

    if result["result"] != "ongoing" {
        match.Status = "finished"
        match.EndTime = time.Now()
        match.Winner = playerID
        if result["result"] == "draw" {
            match.Winner = -1 // Use -1 to indicate a draw
        }
        if err := updateMatchInDatabase(match); err != nil {
            log.Printf("Failed to update match in database: %v", err)
        }
        delete(matches, matchID) // Remove finished game from memory
    }

    return msgBytes, nil
}

func UpdateMatchConfig(broadcastFunc func(int64, []byte), matchID int64, playerID int64, config MatchConfig) error {
	match, exists := matches[matchID]
	if !exists {
		return errors.New("match not found")
	}
	if match.Player1ID != playerID {
		return errors.New("only the match owner can update the configuration")
	}

	match.Config = config

	return broadcastGameInfo(broadcastFunc, matchID)
}

func EndMatch(matchID int64, endTime time.Time, winner int64) error {
	match, err := GetMatch(matchID)
	if err != nil {
		return err
	}

	match.Status = "finished"
	match.EndTime = endTime
	match.Winner = winner

	return updateMatchInDatabase(match)
}