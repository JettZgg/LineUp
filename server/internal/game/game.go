package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/JettZgg/LineUp/internal/db"
	"github.com/JettZgg/LineUp/internal/utils"
	"github.com/JettZgg/LineUp/internal/utils/websocket"
)

type MatchConfig struct {
	BoardWidth  int `json:"boardWidth"`
	BoardHeight int `json:"boardHeight"`
	WinLength   int `json:"winLength"`
}

type Match struct {
	MID              int64       `json:"id"`
	Board            [][]string  `json:"board"`
	Player1ID        int64       `json:"player1Id"`
	Player2ID        int64       `json:"player2Id"`
	Status           string      `json:"status"`
	Config           MatchConfig `json:"config"`
	StartTime        time.Time   `json:"startTime"`
	FirstMovePlayerID int64       `json:"firstMovePlayerId"`
	Moves            []db.Move   `json:"moves"`
}

func (m Match) MarshalJSON() ([]byte, error) {
    type Alias Match
    return json.Marshal(&struct {
        MID string `json:"id"`
        Player1ID string `json:"player1Id"`
        Player2ID string `json:"player2Id"`
        FirstMovePlayerID string `json:"firstMovePlayerId"`
        Alias
    }{
        MID:   strconv.FormatInt(m.MID, 10),
        Player1ID: strconv.FormatInt(m.Player1ID, 10),
        Player2ID: strconv.FormatInt(m.Player2ID, 10),
        FirstMovePlayerID: strconv.FormatInt(m.FirstMovePlayerID, 10),
        Alias: Alias(m),
    })
}

var matches = make(map[int64]*Match)

func CreateMatch(config MatchConfig, playerID int64) (*Match, error) {
	matchID := utils.GenerateMID()
	match := &Match{
		MID:              matchID,
		Board:            makeBoard(config.BoardWidth, config.BoardHeight),
		Player1ID:        playerID,
		Status:           "waiting",
		Config:           config,
		StartTime:        time.Now().UTC(),
		FirstMovePlayerID: playerID,
	}

	dbMatch := &db.Match{
		MID:               match.MID,
		Player1ID:         playerID,
		Status:            match.Status,
		StartTime:         match.StartTime,
		BoardWidth:        config.BoardWidth,
		BoardHeight:       config.BoardHeight,
		WinLength:         config.WinLength,
		FirstMovePlayerID: playerID,
	}

	if err := db.CreateMatch(dbMatch); err != nil {
		return nil, fmt.Errorf("failed to create match in database: %w", err)
	}

	matches[match.MID] = match
	return match, nil
}

func JoinMatch(matchID int64, playerID int64) error {
	match, exists := matches[matchID]
	if !exists {
		return errors.New("match not found")
	}
	if match.Status != "waiting" {
		return errors.New("match is not available for joining")
	}

	// Check if the player is new to the match
	isNewPlayer := match.Player1ID != playerID && match.Player2ID != playerID

	if isNewPlayer {
		if match.Player2ID == 0 {
			match.Player2ID = playerID
			match.Status = "ongoing"

			// Update the match in the database only if a new player joined
			dbMatch := &db.Match{
				MID:       match.MID,
				Player2ID: playerID,
				Status:    match.Status,
			}
			if err := db.UpdateMatch(dbMatch); err != nil {
				log.Printf("Error updating match in database: %v", err)
				return fmt.Errorf("failed to update match in database: %w", err)
			}

			log.Printf("Player %d joined match %d", playerID, matchID)
		} else {
			return errors.New("match is full")
		}
	} else {
		log.Printf("Player %d rejoined match %d", playerID, matchID)
	}

	return nil
}

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

func MakeMove(hub *websocket.Hub, matchID int64, playerID int64, x, y int) (map[string]interface{}, error) {
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
	}

	msgBytes, _ := json.Marshal(updateMsg)
	hub.BroadcastToMatch(matchID, msgBytes)

	if result["result"] != "ongoing" {
		match.Status = "finished"
		endTime := time.Now()
		winner := playerID
		if result["result"] == "draw" {
			winner = -1 // Use -1 to indicate a draw
		}
		if err := updateMatchInDatabase(match, endTime, winner); err != nil {
			log.Printf("Failed to update match in database: %v", err)
		}
		delete(matches, matchID) // Remove finished game from memory
	}

	return result, nil
}

func updateMatchInDatabase(match *Match, endTime time.Time, winner int64) error {
	dbMatch := &db.Match{
		MID:       match.MID,
		Player2ID: match.Player2ID,
		Status:    match.Status,
		EndTime:   endTime,
		Winner:    winner,
		Moves:     match.Moves,
	}
	return db.UpdateMatch(dbMatch)
}

func GetMatchHistory(userID int64, limit int) ([]db.Match, error) {
	return db.GetRecentMatchesByUser(userID, limit)
}

func GetMatchReplay(matchID int64) (*db.Match, error) {
	return db.GetMatchByID(matchID)
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

func checkGameResult(match *Match, x, y int) map[string]interface{} {
	if checkWin(match.Board, x, y, match.Config.WinLength) {
		return map[string]interface{}{"result": "win"}
	}
	if isBoardFull(match.Board) {
		return map[string]interface{}{"result": "draw"}
	}
	return map[string]interface{}{"result": "ongoing"}
}