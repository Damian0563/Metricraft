package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Session struct {
	frontend *websocket.Conn
	worker   *websocket.Conn
	mu       sync.Mutex
}

type Token struct {
	token string
}

func (t Token) GetUser() (User, error) {
	return getUserByToken(t.token)
}

type User struct {
	Mail     string   `json:"mail"`
	AppName  string   `json:"appName"`
	UUID     string   `json:"uuid"`
	Settings settings `json:"settings,omitempty"`
}

// This can be expanded later to include more settings
type settings struct {
	Realtime bool `json:"realtime"`
}
