package main
import (
	"flag"
)

// import "github.com/gorilla/websocket"

var addr = flag.String("addr", ":8080", "http service address")

// func main() {
// 	game := NewGame(NewDdzGameRule())
// 	go game.run()
//
// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		serve(game, w, r)
// 	})
//     fmt.Println("Listening on", *addr)
//     err := http.ListenAndServe(*addr, nil)
//     if err != nil {
//         log.Fatal("ListenAndServe: ", err)
//     }
// }

func main() {
    port := flag.Int("port", 8080, "port to listen on");
    flag.Parse();
}
