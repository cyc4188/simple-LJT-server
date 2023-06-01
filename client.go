package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	conn *websocket.Conn // the websocket connection
	game *Game           // the game the client is in
	send chan []byte
	id   uint
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
    // welcome the player
    c.conn.WriteMessage(websocket.TextMessage, []byte("welcome!"))

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		default:
			return
		}
	}
}

func serve(game *Game, w http.ResponseWriter, r *http.Request) {
	// upgrade http to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// create a new client
	client := &Client{
		conn: conn,
		game: game,
		send: make(chan []byte, 256),
	}
	go client.readMessage()
	go client.sendMessage()
    game.Add <- client
}
