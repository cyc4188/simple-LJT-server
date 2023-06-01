package main

import "net/http"

// import "github.com/gorilla/websocket"

func main() {
	game := NewGame(NewDdzGameRule())
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serve(game, w, r)
	})

	game.run()
}
