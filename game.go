package main

import (
	"errors"
	"log"
	"math/rand"
	"reflect"
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
type GameFail struct { // 出牌失败
	Change
	msg    string
	player *Player
}

type GameStatus struct {
	Change
	score          int
	current_cards  []Card
	current_player *Player
	current_index  int
	Players        []*Player
}

type Action interface {
}

type PlayCards struct {
	Action
	cards  []Card
	player *Player
}
type Pass struct {
	Action
	player *Player
}

type Game struct {
	clients    map[*Client]uuid.UUID // all connected clients and ites id
	gameRule   GameRule
	Mu         sync.RWMutex
	GameState  State // wait, playing, end
	ChangeChan chan Change
	ActionChan chan Action
	GameStatus GameStatus
}

func NewGame(gameRule GameRule) *Game {
	return &Game{
		clients:    make(map[*Client]uuid.UUID),
		gameRule:   gameRule,
		GameState:  Waiting,
		ChangeChan: make(chan Change, 1),
		ActionChan: make(chan Action, 1),
	}
}

func (game *Game) isGameFull() bool {
	return false // TODO: delete
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
		client.Player.client = client
		index++
	}
	game.dealCards()
	game.initGameStatus()
	game.sendChange(game.getStatus()) // 发送游戏状态
	go game.watchAction()
}

// randomly deal cards to players
func (game *Game) dealCards() {
	decks := game.gameRule.generateDeck()

	i := 0
	cardsPerPlayer := game.gameRule.CardPerPlayer()
	for client := range game.clients {
		client.Player.Cards = make([]Card, cardsPerPlayer)
		client.Player.Cards = decks[i : i+cardsPerPlayer]
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
		case Pass:
			game.handlePass(action.(Pass))
		}
		game.sendChange(game.getStatus())
	}
}

func (game *Game) handlePlayCards(play_cards PlayCards) {
	log.Println("handle play cards")
	_, err := game.checkPlayCards(play_cards)
	if err != nil {
		// send error message to client
		game.sendChange(GameFail{
			msg:    err.Error(),
			player: play_cards.player,
		})
		return
	}
	// 更新游戏状态
	game.GameStatus.current_cards = play_cards.cards // 更新当前牌
	play_cards.player.PlayCards(play_cards.cards)    // 玩家出牌
	game.nextPlayer()
	game.GameStatus.current_player.LastPlayed = make([]Card, 0)

	// 广播游戏状态
	game.sendChange(game.getStatus())
}

func (game *Game) handlePass(pass Pass) {
	_, err := game.checkPass(pass)
	if err != nil {
		return
	}
	game.nextPlayer()
	// 判断是否回到出牌方，如果是，则清空当前牌
	if reflect.DeepEqual(game.GameStatus.current_cards, game.GameStatus.current_player.LastPlayed) {
		game.GameStatus.current_cards = make([]Card, 0)
	}
	game.GameStatus.current_player.LastPlayed = make([]Card, 0)

	game.sendChange(game.getStatus())
}

func (game *Game) checkPlayCards(play_cards PlayCards) (bool, error) {
	log.Println("check play cards: ", play_cards)
	// 0. check if it is current player
	if play_cards.player != game.GameStatus.current_player {
		log.Println("not current player")
		return false, errors.New("not current player")
	}
	// 1. check the player has the cards
	if play_cards.player.HasCards(play_cards.cards) == false {
		log.Println("player does not have the cards")
		return false, errors.New("player does not have the cards")
	}
	// 2. check if cards are valid
	if game.gameRule.checkHandsIsValid(play_cards.cards) == false {
		log.Println("invalid hands")
		return false, errors.New("invalid hands")
	}
	// 3. check if cards are bigger than current cards
	if game.gameRule.CompareHands(
		play_cards.cards,
		game.GameStatus.current_cards) != 1 {
		log.Println("invalid hands")
		return false, errors.New("invalid hands")
	}
	return true, nil
}
func (game *Game) checkPass(pass Pass) (bool, error) {
	if pass.player != game.GameStatus.current_player {
		return false, errors.New("not current player")
	}
	return true, nil
}

// move to next player
// will change game status
func (game *Game) nextPlayer() {
	game.GameStatus.current_index = (game.GameStatus.current_index + 1) % len(game.GameStatus.Players)
	game.GameStatus.current_player = game.GameStatus.Players[game.GameStatus.current_index]
}

func (game *Game) sendChange(change Change) {
	game.ChangeChan <- change
}
