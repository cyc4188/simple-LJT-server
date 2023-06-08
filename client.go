package main

import (
	"LJT-server/proto"

	"github.com/google/uuid"
)

type Client struct {
    streamServer proto.Game_StreamServer
    uuid uuid.UUID
    name string
    Player *Player
    game *Game
}

func NewClient(uuid uuid.UUID, name string, game *Game) *Client {
    return &Client{
        uuid: uuid,
        name: name,
        game: game,
    }
}
