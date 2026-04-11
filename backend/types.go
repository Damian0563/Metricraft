package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type signPayload struct {
	Mail    string `json:"mail"`
	Secret  string `json:"secret"`
	AppName string `json:"appName",omitempty`
}

type dashboardInitPayload struct {
	AppName      string `json:"appName"`
	SignedSecret string `json:"signedSecret"`
	Error        string `json:"error"`
}

type Session struct {
	frontend *websocket.Conn
	worker   *websocket.Conn
	mu       sync.Mutex
}
