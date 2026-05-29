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

type existsErrResponse struct {
	Exists bool
	Err    error
	Origin string
}

type Token struct {
	token string
}
type dashboardInitPayload struct {
	AppName      string   `json:"appName"`
	SignedSecret string   `json:"signedSecret"`
	Settings     Settings `json:"settings"`
	Error        string   `json:"error"`
	Urls         []string `json:"urls"`
}

// This can be expanded later to include more settings
type Settings struct {
	Realtime  bool            `json:"realtime"`
	Enabled   map[string]bool `json:"enabled"`
	Retention int             `json:"retention"`
}

type Metric struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
