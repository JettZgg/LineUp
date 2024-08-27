// File: internal/websocket/client.go
package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/JettZgg/LineUp/internal/match"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	roomID int64
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		switch msg["type"] {
		case "getMatchInfo":
			matchID := int64(msg["matchId"].(float64))
			matchInfo, err := match.GetMatchInfo(matchID)
			if err != nil {
				log.Printf("Error getting match info: %v", err)
				continue
			}
			response, _ := json.Marshal(matchInfo)
			c.send <- response
		case "startMatch":
			matchID := int64(msg["matchId"].(float64))
			userID := int64(msg["userId"].(float64))
			if err := match.StartMatch(c.hub.BroadcastToMatch, matchID, userID); err != nil {
				log.Printf("Error starting match: %v", err)
				continue
			}
		case "makeMove":
			matchID := int64(msg["matchId"].(float64))
			userID := int64(msg["userId"].(float64))
			move := msg["move"].(string)
			result, err := match.MakeMove(c.hub.BroadcastToMatch, matchID, userID, move)
			if err != nil {
				log.Printf("Error making move: %v", err)
				errorMsg := map[string]interface{}{
					"type":  "error",
					"error": err.Error(),
				}
				errorJSON, _ := json.Marshal(errorMsg)
				c.send <- errorJSON
			} else {
				c.hub.BroadcastToMatch(matchID, result)
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
