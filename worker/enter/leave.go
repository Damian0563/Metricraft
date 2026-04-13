package enter

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

func Leave(payload Payload) error {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws/workers", nil)
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
	err = Connect()
	fmt.Println(err)
	return err
}
