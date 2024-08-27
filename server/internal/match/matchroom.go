package match

import (
	"errors"
	"fmt"
	"log"

	"github.com/JettZgg/LineUp/internal/db"
	"golang.org/x/exp/rand"
)

func StartMatch(broadcastFunc func(int64, []byte), matchID int64, playerID int64) error {
	match, exists := matches[matchID]
	if !exists {
		return errors.New("match not found")
	}
	if match.Player1ID != playerID {
		return errors.New("only the match owner can start the match")
	}
	if match.Player2ID == 0 {
		return errors.New("waiting for second player to join")
	}

	match.FirstMovePlayerID = determineFirstPlayer(match.Player1ID, match.Player2ID)

	// Update the match in the database
	dbMatch := &db.Match{
		MID:               match.MID,
		FirstMovePlayerID: match.FirstMovePlayerID,
	}
	if err := db.UpdateMatch(dbMatch); err != nil {
		log.Printf("Error updating match in database: %v", err)
		return fmt.Errorf("failed to update match in database: %w", err)
	}

	return broadcastMatchInfo(broadcastFunc, matchID)
}

func determineFirstPlayer(player1ID, player2ID int64) int64 {
	if rand.Intn(2) == 0 {
		return player1ID
	}
	return player2ID
}
