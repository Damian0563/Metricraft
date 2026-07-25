package types

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
