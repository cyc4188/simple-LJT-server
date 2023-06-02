package main

type GameState struct {
    currentClient *Client // the client whose turn it is
    game *Game
}
