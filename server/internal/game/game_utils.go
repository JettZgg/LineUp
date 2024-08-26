package game

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/JettZgg/LineUp/internal/db"
)

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

func GetGameInfo(matchID int64) (map[string]interface{}, error) {
	match, err := GetMatch(matchID)
	if err != nil {
		log.Printf("Error getting match %d: %v", matchID, err)
		return nil, err
	}

	log.Printf("Match data: %+v", match)

	players := []map[string]interface{}{}

	if match.Player1ID != 0 {
		player1, err := db.GetUserByID(match.Player1ID)
		if err != nil {
			log.Printf("Error getting player1 (ID: %d): %v", match.Player1ID, err)
			return nil, err
		}
		players = append(players, map[string]interface{}{
			"id":       player1.UID,
			"username": player1.Username,
			"ready":    match.Player1Ready,
		})
	}

	if match.Player2ID != 0 {
		player2, err := db.GetUserByID(match.Player2ID)
		if err != nil {
			log.Printf("Error getting player2 (ID: %d): %v", match.Player2ID, err)
			return nil, err
		}
		players = append(players, map[string]interface{}{
			"id":       player2.UID,
			"username": player2.Username,
			"ready":    match.Player2Ready,
		})
	}

	log.Printf("GetGameInfo for match %d: players=%v, config=%v", matchID, players, match.Config)

	return map[string]interface{}{
		"type":    "gameInfo",
		"matchId": matchID,
		"players": players,
		"config":  match.Config,
		"status":  match.Status,
		"board":   match.Board,
	}, nil
}

func broadcastGameInfo(broadcastFunc func(int64, []byte), matchID int64) error {
	gameInfo, err := GetGameInfo(matchID)
	if err != nil {
		return fmt.Errorf("failed to get game info: %w", err)
	}

	msgBytes, err := json.Marshal(gameInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal game info: %w", err)
	}

	if string(msgBytes) == "null" || string(msgBytes) == "" {
		log.Printf("Warning: Attempting to broadcast empty game info for match %d", matchID)
		return nil
	}

	log.Printf("Broadcasting game info for match %d: %s", matchID, string(msgBytes))
	broadcastFunc(matchID, msgBytes)

	return nil
}

func updateMatchInDatabase(match *Match) error {
	dbMatch := &db.Match{
		MID:       match.MID,
		Player2ID: match.Player2ID,
		Status:    match.Status,
		EndTime:   match.EndTime,
		Winner:    match.Winner,
		Moves:     match.Moves,
	}
	return db.UpdateMatch(dbMatch)
}