package api

import (
	"net/http"

	"github.com/JettZgg/LineUp/internal/middleware"
	"github.com/JettZgg/LineUp/internal/utils/websocket"
	"github.com/gorilla/mux"
)

func SetupRoutes(hub *websocket.Hub) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")

	r.HandleFunc("/api/create-match", middleware.AuthMiddleware(CreateMatchHandler)).Methods("POST")
	r.HandleFunc("/api/join-match/{matchID}", middleware.AuthMiddleware(JoinMatchHandler)).Methods("POST")
	r.HandleFunc("/api/make-move/{matchID}", middleware.AuthMiddleware(MakeMoveHandler)).Methods("POST")

	r.HandleFunc("/ws/{matchID}", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	return r
}
