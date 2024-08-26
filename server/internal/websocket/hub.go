package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients    map[*Client]bool
	Rooms      map[int64]*Room
	Register   chan *Client
	Unregister chan *Client
	Matches    map[int64]map[*Client]bool
}

type Room struct {
	ID      int64
	Clients map[*Client]bool
	Send    chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Rooms:      make(map[int64]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Matches:    make(map[int64]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			room, ok := h.Rooms[client.roomID]
			if !ok {
				room = &Room{
					ID:      client.roomID,
					Clients: make(map[*Client]bool),
					Send:    make(chan []byte),
				}
				h.Rooms[client.roomID] = room
				go room.run()
			}
			room.Clients[client] = true
			h.Matches[client.roomID] = room.Clients
		case client := <-h.Unregister:
			if room, ok := h.Rooms[client.roomID]; ok {
				if _, ok := room.Clients[client]; ok {
					delete(room.Clients, client)
					close(client.send)
					if len(room.Clients) == 0 {
						close(room.Send)
						delete(h.Rooms, client.roomID)
						delete(h.Matches, client.roomID)
					}
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

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, roomID int64) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), roomID: roomID}
	client.hub.Register <- client

	go client.writePump()
	go client.readPump()
}

func (h *Hub) BroadcastToMatch(matchID int64, message []byte) {
	if clients, ok := h.Matches[matchID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.Clients, client)
				delete(h.Matches[matchID], client)
			}
		}
	}
}

func (r *Room) run() {
	for {
		select {
		case message := <-r.Send:
			for client := range r.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.Clients, client)
				}
			}
		}
	}
}