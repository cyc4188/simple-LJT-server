package main

import (
	"sync"

	"github.com/google/uuid"
)

type State int

const (
    Waiting State = iota
    Playing
    End
)

type Change interface {
}

type GameStart struct {
    Change
}

type GameEnd struct {
    Change
}

type GameStatus struct {
    Change
    // TODO
}

type Game struct {
	clients  map[*Client]uuid.UUID // all connected clients and ites id
	gameRule GameRule
    score int
    Mu sync.RWMutex
    GameState State // wait, playing, end
    ChangeChan chan Change
}

func NewGame(gameRule GameRule) *Game {
	return &Game{
		clients:  make(map[*Client]uuid.UUID),
		gameRule: gameRule,
        GameState: Waiting,
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
// return true if player count is equal to game rule
func (game *Game) checkPlayerCount() bool {
	return len(game.clients) == game.gameRule.PlayerCount()
}

func (game *Game) startGame() {
    game.GameState = Playing
    game.score = 0
    for client := range game.clients {
        client.Player = NewPlayer(client)
    }
    game.dealCards()
    go game.watchAction()
}

// 监听玩家动作，如出牌，跳过等
func (game *Game) watchAction() {
    // TODO
}

func (game *Game) SendChange(change Change) {
    // TODO
}

func (game *Game) dealCards() {
    // TODO
    decks := game.gameRule.generateDeck()
    
    i := 0
    cardsPerPlayer := game.gameRule.CardPerPlayer()
    for client := range game.clients {
        client.Player.cards = make([]Card, cardsPerPlayer)
        client.Player.cards = decks[i:i+cardsPerPlayer]
        i += cardsPerPlayer
    }
    // TODO
    // send status to server
}

