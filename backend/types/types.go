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

type ExistsErrResponse struct {
	Exists bool
	Err    error
	Origin string
	Owner  string
}

type DashboardInitPayload struct {
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

type PendingUsers struct {
	Mail string `json:"mail"`
}

type AllowedUsers struct {
	Mail     string `json:"mail"`
	Initials string `json:"initials"`
}
type SendRecoveryUser struct {
	Mail string `json:"mail"`
	Id   string `json:"id:omitempty"`
}
