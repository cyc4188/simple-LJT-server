package main

import (
	"LJT-server/proto"
	"math/rand"
	"sync"
	"time"
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
    score int
    current_cards []Card
    current_player *Player
    Players []*Player
}

type Action interface {
}

type PlayCards struct {
    Action
    cards []Card
    player *Player
}
type Pass struct {
    Action
    player *Player
}

func (gamestatus *GameStatus) ToProto(player *Player) *proto.StreamResponse {
    players := make([]*proto.Player, 0)
    for _, player := range gamestatus.Players {
        players = append(players, player.ToProto())
    }
    current_cards := make([]*proto.Card, 0)
    for _, card := range gamestatus.current_cards {
        current_cards = append(current_cards, card.ToProto())
    }
    return &proto.StreamResponse{
        Response: &proto.StreamResponse_Continue{
            Continue: &proto.Continue{
                Score: int32(gamestatus.score),
                CurrentCards: current_cards,
                CurrentPlayer: gamestatus.current_player.ToProto(),
                Players: players,
            },
        },
    }
}

type Game struct {
	clients  map[*Client]uuid.UUID // all connected clients and ites id
	gameRule GameRule
    Mu sync.RWMutex
    GameState State // wait, playing, end
    ChangeChan chan Change
    ActionChan chan Action
    GameStatus GameStatus
}

func NewGame(gameRule GameRule) *Game {
	return &Game{
		clients:  make(map[*Client]uuid.UUID),
		gameRule: gameRule,
        GameState: Waiting,
        ChangeChan: make(chan Change, 1),
        ActionChan: make(chan Action, 1),
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

// return game status
func (game *Game) getStatus() GameStatus {
    return game.GameStatus
}

func (game *Game) initGameStatus() {
    game.GameStatus.score = 0
    game.GameStatus.current_cards = make([]Card, 0)

    game.GameStatus.Players = make([]*Player, 0)
    for client := range game.clients {
        game.GameStatus.Players = append(game.GameStatus.Players, client.Player)
    }

    game.GameStatus.current_player = game.GameStatus.Players[rand.Intn(len(game.GameStatus.Players))]
}

// start game
// deal cards
// watch for action
func (game *Game) startGame() {
    // sleep for 3 seconds to wait for other players
    time.Sleep(3 * time.Second)

    game.GameState = Playing
    index := 0
    for client := range game.clients {
        client.Player = NewPlayer(client)
        client.Player.index = uint(index)
        index++
    }
    game.dealCards()
    game.initGameStatus()
    game.SendChange(game.getStatus()) // 发送游戏状态
    go game.watchAction()
}

// randomly deal cards to players
func (game *Game) dealCards() {
    decks := game.gameRule.generateDeck()
    
    i := 0
    cardsPerPlayer := game.gameRule.CardPerPlayer()
    for client := range game.clients {
        client.Player.Cards = make([]Card, cardsPerPlayer)
        client.Player.Cards = decks[i:i+cardsPerPlayer]
        i += cardsPerPlayer
    }
}

// 监听玩家动作，如出牌，跳过等
func (game *Game) watchAction() {
    for {
        action := <-game.ActionChan
        switch action.(type) {
        case PlayCards:
            game.handlePlayCards(action.(PlayCards))
        }
        game.SendChange(game.getStatus())
    }
}

func (game *Game) handlePlayCards(play_cards PlayCards) {
    // 1. check if cards are valid 

    // 2. check if cards are bigger than current cards
}

func (game *Game) handlePass(pass Pass) {
}


func (game *Game) SendChange(change Change) {
    game.ChangeChan <- change
}
