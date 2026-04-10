package enter

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

func Leave(payload Payload) error {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000/visitors", nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}
