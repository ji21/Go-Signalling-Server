package main

import (
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type WSMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
	User User   `json:"User"`
}

type User struct {
	ID   int
	conn *websocket.Conn
}

type Negotiation struct {
	Type string                    `json"type"`
	SDP  webrtc.SessionDescription `json:"SDP"`
}
