package enter

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"os"
	"worker/db"
	"worker/types"
)

func Leave(payload types.Payload) error {
	host := os.Getenv("ws")
	conn, _, err := websocket.DefaultDialer.Dial(host+":8080/ws/workers", nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}
	return db.Insert(payload)
}
