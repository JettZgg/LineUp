package game

import (
	"errors"
	"fmt"
)

type MatchConfig struct {
	BoardWidth  int `json:"boardWidth"`
	BoardHeight int `json:"boardHeight"`
	WinLength   int `json:"winLength"`
}

type Match struct {
	ID        string      `json:"id"`
	Board     [][]string  `json:"board"`
	Player1ID string      `json:"player1Id"`
	Player2ID string      `json:"player2Id"`
	Status    string      `json:"status"`
	Config    MatchConfig `json:"config"`
}

var matches = make(map[string]*Match)

func CreateMatch(config MatchConfig) (*Match, error) {
	match := &Match{
		ID:     generateMatchID(),
		Board:  makeBoard(config.BoardWidth, config.BoardHeight),
		Status: "waiting",
		Config: config,
	}
	matches[match.ID] = match
	return match, nil
}

func JoinMatch(matchID, playerID string) error {
	match, exists := matches[matchID]
	if !exists {
		return errors.New("match not found")
	}
	if match.Status != "waiting" {
		return errors.New("match is not available for joining")
	}
	if match.Player1ID == "" {
		match.Player1ID = playerID
	} else {
		match.Player2ID = playerID
		match.Status = "ongoing"
	}
	return nil
}

func MakeMove(matchID, playerID string, x, y int) (map[string]interface{}, error) {
	match, exists := matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
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

	if checkWin(match.Board, x, y, match.Config.WinLength) {
		match.Status = "finished"
		return map[string]interface{}{"result": "win", "winner": playerID}, nil
	}
	if isBoardFull(match.Board) {
		match.Status = "finished"
		return map[string]interface{}{"result": "draw"}, nil
	}
	return map[string]interface{}{"result": "ongoing"}, nil
}

func makeBoard(width, height int) [][]string {
	board := make([][]string, height)
	for i := range board {
		board[i] = make([]string, width)
	}
	return board
}

func checkWin(board [][]string, x, y, winLength int) bool {
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	symbol := board[y][x]

	for _, dir := range directions {
		count := 1
		for i := 1; i < winLength; i++ {
			nx, ny := x+dir[0]*i, y+dir[1]*i
			if nx < 0 || ny < 0 || nx >= len(board[0]) || ny >= len(board) || board[ny][nx] != symbol {
				break
			}
			count++
		}
		for i := 1; i < winLength; i++ {
			nx, ny := x-dir[0]*i, y-dir[1]*i
			if nx < 0 || ny < 0 || nx >= len(board[0]) || ny >= len(board) || board[ny][nx] != symbol {
				break
			}
			count++
		}
		if count >= winLength {
			return true
		}
	}
	return false
}

func isBoardFull(board [][]string) bool {
	for _, row := range board {
		for _, cell := range row {
			if cell == "" {
				return false
			}
		}
	}
	return true
}

func generateMatchID() string {
	// Implement a proper unique ID generation mechanism
	return fmt.Sprintf("match_%d", len(matches)+1)
}
