package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var isStreaming = false

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func reader(conn *websocket.Conn) {
	for {
		wsMessage := WSMessage{}
		if err := conn.ReadJSON(&wsMessage); err != nil {
			log.Fatal(err)
		}
		if wsMessage.Type == "offer" {
			log.Println("offer")
		} else if wsMessage.Type == "answer" {
			log.Println("answer")
		} else if wsMessage.Type == "join" {
			log.Println("join")
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	reader(conn)
}

func setRoutes() {
	http.HandleFunc("/ws", wsHandler)
}

func main() {
	setRoutes()
	log.Println("Starting server on port 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
