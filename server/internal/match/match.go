package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/JettZgg/LineUp/internal/db"
)

type Match struct {
	MID               int64  `json:"id"`
	Player1ID         int64  `json:"player1Id"`
	Player2ID         int64  `json:"player2Id"`
	Winner            int64  `json:"winner"`
	FirstMovePlayerID int64  `json:"firstMovePlayerId"`
	Moves             string `json:"moves"`
	Date              string `json:"date"`
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
			Winner:            dbMatch.Winner,
			FirstMovePlayerID: dbMatch.FirstMovePlayerID,
			Moves:             dbMatch.Moves,
			Date:              dbMatch.Date,
		}
		matches[matchID] = match
	}
	return match, nil
}

func MakeMove(broadcastFunc func(int64, []byte), matchID int64, playerID int64, move string) ([]byte, error) {
	match, err := GetMatch(matchID)
	if err != nil {
		return nil, err
	}
	if match.Winner != 0 {
		return nil, errors.New("match is already finished")
	}
	if match.Player1ID != playerID && match.Player2ID != playerID {
		return nil, errors.New("invalid player")
	}

	match.Moves += move

	// Check for win condition (implement this function based on Gomoku rules)
	if checkWin(match.Moves) {
		match.Winner = playerID
	}

	// Create a message to broadcast
	updateMsg := map[string]interface{}{
		"type":    "moveUpdate",
		"matchID": matchID,
		"moves":   match.Moves,
		"lastMove": map[string]interface{}{
			"playerID": playerID,
			"move":     move,
		},
	}

	if match.Winner != 0 {
		updateMsg["winner"] = match.Winner
	}

	msgBytes, err := json.Marshal(updateMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal move update: %w", err)
	}

	if err := updateMatchInDatabase(match); err != nil {
		log.Printf("Failed to update match in database: %v", err)
	}

	return msgBytes, nil
}

func EndMatch(matchID int64, winner int64) error {
	match, err := GetMatch(matchID)
	if err != nil {
		return err
	}

	match.Winner = winner

	return updateMatchInDatabase(match)
}

func checkWin(moves string) bool {
	board := make([][]rune, 15)
	for i := range board {
		board[i] = make([]rune, 15)
	}

	// Fill the board with moves
	for i := 0; i < len(moves); i += 3 {
		if i+2 >= len(moves) {
			break
		}
		col := int(moves[i] - 'a')
		row, err := strconv.Atoi(moves[i+1 : i+3])
		if err != nil {
			log.Printf("Error parsing row number: %v", err)
			continue
		}
		row = 15 - row // Invert the row number
		if row < 0 || row >= 15 || col < 0 || col >= 15 {
			log.Printf("Invalid move: %s", moves[i:i+3])
			continue
		}
		if i/3%2 == 0 {
			board[row][col] = 'X'
		} else {
			board[row][col] = 'O'
		}
	}

	// Check for 5 in a row
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if board[i][j] == 0 {
				continue
			}
			// Check horizontal
			if j <= 10 && board[i][j] == board[i][j+1] && board[i][j] == board[i][j+2] && board[i][j] == board[i][j+3] && board[i][j] == board[i][j+4] {
				return true
			}
			// Check vertical
			if i <= 10 && board[i][j] == board[i+1][j] && board[i][j] == board[i+2][j] && board[i][j] == board[i+3][j] && board[i][j] == board[i+4][j] {
				return true
			}
			// Check diagonal (top-left to bottom-right)
			if i <= 10 && j <= 10 && board[i][j] == board[i+1][j+1] && board[i][j] == board[i+2][j+2] && board[i][j] == board[i+3][j+3] && board[i][j] == board[i+4][j+4] {
				return true
			}
			// Check diagonal (top-right to bottom-left)
			if i <= 10 && j >= 4 && board[i][j] == board[i+1][j-1] && board[i][j] == board[i+2][j-2] && board[i][j] == board[i+3][j-3] && board[i][j] == board[i+4][j-4] {
				return true
			}
		}
	}
	return false
}