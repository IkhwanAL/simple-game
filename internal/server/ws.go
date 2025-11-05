package server

import (
	"context"
	"sync"

	"github.com/coder/websocket"
)

type WebSocketHub struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]bool
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (hub *WebSocketHub) AddConn(conn *websocket.Conn) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	hub.conns[conn] = true
}

func (hub *WebSocketHub) RemoveConn(conn *websocket.Conn) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	delete(hub.conns, conn)
}

func (hub *WebSocketHub) Broadcast(ctx context.Context, msg []byte) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	for conn := range hub.conns {
		err := conn.Write(ctx, websocket.MessageText, msg)
		if err != nil {
			conn.Close(websocket.StatusNormalClosure, "disconnect")
			delete(hub.conns, conn)
		}
	}
}
