package types

import "github.com/gorilla/websocket"

type Client struct {
	Connection *websocket.Conn
	Type       string
}
