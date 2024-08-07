package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	matches    map[int64]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		matches:    make(map[int64]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if _, ok := h.matches[client.matchID]; !ok {
				h.matches[client.matchID] = make(map[*Client]bool)
			}
			h.matches[client.matchID][client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.matches[client.matchID], client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// Parse the message to get the matchID
			var msg struct {
				MatchID int64 `json:"matchID"`
				// other fields...
			}
			json.Unmarshal(message, &msg)

			for client := range h.matches[msg.MatchID] {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					delete(h.matches[msg.MatchID], client)
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust this for production
	},
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, matchID int64) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), matchID: matchID}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (h *Hub) BroadcastToMatch(matchID int64, message []byte) {
	if clients, ok := h.matches[matchID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(h.matches[matchID], client)
			}
		}
	}
}
