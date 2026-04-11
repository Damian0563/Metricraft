package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var session = Session{
	frontend: nil,
	worker:   nil,
	mu:       sync.Mutex{},
}

func wsHandlerFrontend(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	session.mu.Lock()
	session.frontend = conn
	session.mu.Unlock()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			session.mu.Lock()
			if session.frontend == conn {
				session.frontend = nil
			}
			session.mu.Unlock()
			break
		}
	}
}

func wsHandlerWorker(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	session.mu.Lock()
	session.worker = conn
	session.mu.Unlock()
	for {
		_, message, err := conn.ReadMessage()
		session.mu.Lock()
		frontend := session.frontend
		session.mu.Unlock()
		if err != nil {
			session.mu.Lock()
			if session.worker == conn {
				session.worker = nil
			}
			session.mu.Unlock()
		}
		if frontend == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		err = frontend.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			session.mu.Lock()
			if session.frontend == frontend {
				session.frontend = nil
			}
			session.mu.Unlock()
		}
	}
}

func StartWebSocketServer(router *http.ServeMux) {
	router.HandleFunc("/ws/visitors", wsHandlerFrontend)
	router.HandleFunc("/ws/workers", wsHandlerWorker)
}
