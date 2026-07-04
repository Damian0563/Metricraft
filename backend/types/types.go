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

type Worker struct {
	Url          string            `json:"url"`
	PollInterval int               `json:"pollInterval"`
	Headers      map[string]string `json:"headers,omitempty"`
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

type Settings struct {
	Realtime  bool                     `json:"realtime"`
	Enabled   map[string]EnabledMetric `json:"enabled"`
	Retention int                      `json:"retention"`
}

type EnabledMetric struct {
	Enabled   bool   `json:"enabled"`
	Timeframe string `json:"timeframe"`
}

type Metric struct {
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	Timeframe string `json:"timeframe"`
}

type PendingUsers struct {
	Mail string `json:"mail"`
}

type AllowedUsers struct {
	Mail     string `json:"mail"`
	Initials string `json:"initials"`
	Status   bool   `json:"status"`
}
type SendRecoveryUser struct {
	Mail string `json:"mail"`
	Id   string `json:"id,omitempty"`
}

type RecoveryCheckPayload struct {
	Password string `json:"password"`
}
type Invite struct {
	Invitees []string `json:"invitees"`
}
