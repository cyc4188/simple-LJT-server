package main

import (
	"log"

	"github.com/gorilla/websocket"
)

const (
    maxMessageSize = 512
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

type Client struct {
    conn *websocket.Conn // the websocket connection
    game *Game          // the game the client is in
    send chan []byte
    id uint
}

// read from client
func (c *Client) readMessage() {
    c.conn.SetReadLimit(maxMessageSize)
    
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            log.Printf("read error: %v", err)
            break
        }
        log.Printf("id: %v recv: %s", c.id, message)
    }
}

// send to the client
func (c *Client) sendMessage() {
    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            c.conn.WriteMessage(websocket.TextMessage, message)
        }
    }
}

func (c *Client) serve() {
    go c.readMessage()
    go c.sendMessage()
}
