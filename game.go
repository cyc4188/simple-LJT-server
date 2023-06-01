package main

type Game struct {
    clients map[*Client]uint // all connected clients and ites id
    gameRule GameRule
}

// check player count
func (game *Game) checkPlayerCount() bool {
    return len(game.clients) == game.gameRule.PlayerCount()
}

// deal cards
func (game *Game) dealCards() {
    // deck := game.gameRule.generateDeck()           
}
