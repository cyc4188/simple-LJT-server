package main

import (
	"github.com/google/uuid"
)

type Client struct {
    uuid uuid.UUID
    name string
    game *Game
}

func NewClient(uuid uuid.UUID, name string, game *Game) *Client {
    return &Client{
        uuid: uuid,
        name: name,
        game: game,
    }
}
