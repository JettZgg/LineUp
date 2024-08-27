// File: internal/websocket/client.go
package websocket

import (
	"encoding/json"
	"log"
	"strconv"
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
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		switch msg["type"] {
		case "getMatchInfo":
			var matchIDInt int64
			switch v := msg["matchId"].(type) {
			case string:
				matchIDInt, err = strconv.ParseInt(v, 10, 64)
				if err != nil {
					log.Printf("Error parsing matchId: %v", err)
					continue
				}
			case float64:
				matchIDInt = int64(v)
			default:
				log.Printf("Invalid matchId type")
				continue
			}
			matchInfo, err := match.GetMatchInfo(matchIDInt)
			if err != nil {
				log.Printf("Error getting match info: %v", err)
				continue
			}
			response, _ := json.Marshal(matchInfo)
			c.send <- response
		case "startMatch":
			matchID := getInt64(msg["matchId"])
			userID := getInt64(msg["userId"])
			if err := match.StartMatch(c.hub.BroadcastToMatch, matchID, userID); err != nil {
				log.Printf("Error starting match: %v", err)
				continue
			}
		case "makeMove":
			matchID := getInt64(msg["matchId"])
			userID := getInt64(msg["userId"])
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

func getInt64(v interface{}) int64 {
	switch i := v.(type) {
	case float64:
		return int64(i)
	case string:
		j, _ := strconv.ParseInt(i, 10, 64)
		return j
	default:
		return 0
	}
}
