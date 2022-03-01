package main

import "github.com/gorilla/websocket"

type WSMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type User struct {
	ID   int
	conn *websocket.Conn
}
