package enter

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"os"
)

func Leave(payload Payload) error {
	port := os.Getenv("PORT")
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:"+port+"/leave", nil)
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
