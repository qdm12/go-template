package websocket

import (
	"golang.org/x/net/websocket"
)

func New() *websocket.Server {
	return &websocket.Server{
		Config: websocket.Config{
			Origin: nil,
		},
	}
}
