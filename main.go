package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var viewers = make(map[User]int)

var streamer = User{ID: -1, conn: nil}

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
			onOffer(&wsMessage)
		} else if wsMessage.Type == "answer" {
			onAnswer(&wsMessage)
		} else if wsMessage.Type == "candidate" {
			onCandidate(&wsMessage)
		} else if wsMessage.Type == "view" {
			onView(&wsMessage)
		}
	}
}

func onView(wsMessage *WSMessage) {
	user := wsMessage.User
	exist := viewers[user]
	if exist == 0 {
		viewers[user] = 0
	} else {
		viewers[user] = 1
	}
}

func onCandidate(wsMessage *WSMessage) {
	if wsMessage.User == streamer {
		for user := range viewers {
			user.conn.WriteJSON(&wsMessage)
		}
	} else {
		streamer.conn.WriteJSON(&wsMessage)
	}
}

func onAnswer(wsMessage *WSMessage) {
	var negotiation Negotiation
	_ = json.Unmarshal([]byte(wsMessage.Data), &negotiation)
	if streamer.ID != -1 {
		streamer.conn.WriteJSON(&negotiation)
	} else {
		log.Println("No streamer to send answer to.")
	}
}

func onOffer(wsMessage *WSMessage) {
	var negotiation Negotiation
	_ = json.Unmarshal([]byte(wsMessage.Data), &negotiation)
	if streamer.ID == -1 {
		streamer = wsMessage.User
	}

	for user := range viewers {
		user.conn.WriteJSON(&negotiation)
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
