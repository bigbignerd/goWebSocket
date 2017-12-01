package main

import (
	"github.com/gorilla/websocket"
	// "log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"token"},
	CheckOrigin:  func(r *http.Request) bool { return true },
}

func main() {

	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		onlineUser(hub, w, r)
	})
	http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(hub, w, r)
	})
	http.ListenAndServe(":3000", nil)
}
