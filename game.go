package main

import (
	"sync"

	"github.com/google/uuid"
)

type Game struct {
	clients  map[*Client]uuid.UUID // all connected clients and ites id
	gameRule GameRule
    Mu sync.RWMutex
}

func NewGame(gameRule GameRule) *Game {
	return &Game{
		clients:  make(map[*Client]uuid.UUID),
		gameRule: gameRule,
	}
}


func (game *Game) isGameFull() bool {
    return game.clientCount() >= game.gameRule.PlayerCount()
}

func (game *Game) clientCount() int {
    return len(game.clients)
}

func (game *Game) addClient(client *Client) {
	game.clients[client] = client.uuid
}

func (game *Game) removeClient(client *Client) {
	delete(game.clients, client)
}

// check player count
func (game *Game) checkPlayerCount() bool {
	return len(game.clients) == game.gameRule.PlayerCount()
}

