package main

type Game struct {
	clients  map[*Client]uint // all connected clients and ites id
	gameRule GameRule
    Add chan *Client
    Remove chan *Client
}

func NewGame(gameRule GameRule) *Game {
	return &Game{
		clients:  make(map[*Client]uint),
		gameRule: gameRule,
	}
}

func (game *Game) addClient(client *Client) {
	game.clients[client] = 0
}

func (game *Game) removeClient(client *Client) {
	delete(game.clients, client)
}

// check player count
func (game *Game) checkPlayerCount() bool {
	return len(game.clients) == game.gameRule.PlayerCount()
}

// deal cards
func (game *Game) dealCards() {
	// deck := game.gameRule.generateDeck()
}

func (game *Game) run() {
    for {
        select {
        case client := <-game.Add:
            game.addClient(client)
        case client := <-game.Remove:
            game.removeClient(client)
        }
    }
}

func (game *Game) gameLoop() {
    for {
        if game.checkPlayerCount() {
            game.dealCards()
        }
    }
}
