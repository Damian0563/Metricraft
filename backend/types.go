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

type User struct {
	Mail     string   `json:"mail"`
	AppName  string   `json:"appName"`
	UUID     string   `json:"uuid"`
	Settings Settings `json:"settings,omitempty"`
}

type Token struct {
	token string
}

// This can be expanded later to include more settings
type Settings struct {
	Realtime bool   `json:"realtime"`
	Enabled  string `json:"enabled"`
}
