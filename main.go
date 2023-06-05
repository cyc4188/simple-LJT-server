package main

import (
	"LJT-server/proto"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

// import "github.com/gorilla/websocket"

var addr = flag.String("addr", ":8080", "http service address")

func main() {
    port := flag.Int("port", 8080, "port to listen on");
    flag.Parse();

    log.Printf("listen on port %d", *port)
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    if err != nil {
        log.Fatalf("failed to listen on %d", *port)
    }

    game := NewGame(NewDdzGameRule())
    
    var s = grpc.NewServer();
    var server = NewGameServer(game)
    proto.RegisterGameServer(s, server)
    
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve on %d", *port)
    }
}
