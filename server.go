package main

import (
	"LJT-server/proto"
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GameServer struct {
    proto.UnimplementedGameServer;
    game *Game; // TODO: support for multiple games match
    clients map[uuid.UUID]*Client;
    mu sync.RWMutex;
}

func NewGameServer(game *Game) *GameServer {
    server :=  &GameServer{
        game: game,
        clients: make(map[uuid.UUID]*Client),
    }
    go server.watchChange()
    return server
}

// 处理连接请求
// refer: proto ConnectRequest
func (s *GameServer) Connecting(ctx context.Context, req *proto.ConnectRequest) (*proto.ConnectResponse, error) {
    if (s.game.isGameFull()) {
        return nil, errors.New("server is full")
    }  
    playerId, err := uuid.Parse(req.Id)
    if err != nil {
        return nil, err
    }
    
    // check if player already connected
    s.mu.RLock()
    if _, ok := s.clients[playerId]; ok {
        s.mu.RUnlock()
        return nil, errors.New("player already connected")
    }
    s.mu.RUnlock()

    newClient := NewClient(playerId, req.Name, s.game)

    log.Printf("player %s connected", req.Name)
    // add client to the server 
    s.mu.Lock()
    s.clients[playerId] = newClient
    s.mu.Unlock()

    // add client to the game 
    s.game.Mu.Lock()
    s.game.addClient(newClient)
    s.game.Mu.Unlock()
    if s.game.checkPlayerCount() {
        s.game.startGame()
    }

    return &proto.ConnectResponse{
        Token: "test",
        Players: []*proto.Player{},
    }, nil
}

func (s *GameServer) Stream(stream proto.Game_StreamServer) error {
    ctx := stream.Context()

    // 获取第一个请求
    req, err := stream.Recv()
    if err != nil {
        log.Printf("receive error: %v", err)
        return nil
    }
    log.Printf("receive: %v", req.PlayCards) 
    playerId, err := uuid.Parse(req.PlayCards.Player.Id)
    // 根据请求的uuid得到对应的streamServer
    s.clients[playerId].streamServer = stream

    go func() {
        for {
            req, err := stream.Recv()
            if err != nil {
                log.Printf("receive error: %v", err)
                return 
            }
            log.Printf("receive: %v", req.PlayCards) 
            // push to action chan 
            s.handleRequest(req)
        }
    }()
    
    select {
    case <-ctx.Done():
        log.Printf("client disconnected")
        // TODO: remove client from game
        return nil
    }
    return nil
}

func (s *GameServer) handleRequest(req *proto.StreamRequest) {
    // TODO
}

func (s *GameServer) watchChange() {
    for {
        change := <-s.game.ChangeChan
        switch change.(type) {
        case GameStatus:
            change := change.(GameStatus)
            s.broadcast(change)
        }
    }
}


// broadcast game status to all clients
// and its own cards
func (s *GameServer) broadcast(status GameStatus) {
    // sleep for 3 seconds to wait for other players
    time.Sleep(3 * time.Second)

    log.Printf("broadcast game status")

    players := make([]*proto.Player, 0)
    for _, player := range status.Players {
        players = append(players, player.ToProto())
    }
    current_cards := make([]*proto.Card, 0)
    for _, card := range status.current_cards {
        current_cards = append(current_cards, card.ToProto())
    }

    for _, client := range s.clients {
        if client.streamServer == nil {
            log.Println("client ", client.uuid, " stream server is nil")
            continue
        }
        cards := make([]*proto.Card, 0)
        for _, card := range client.Player.Cards {
            cards = append(cards, card.ToProto())
        }

        msg := &proto.StreamResponse{  
            Response: &proto.StreamResponse_Continue{
                Continue: &proto.Continue{
                    Score: int32(status.score),
                    Players: players,
                    CurrentCards: current_cards, 
                    CurrentPlayer: status.current_player.ToProto(),
                    Cards: cards,
                },
            },
        }
        client.streamServer.Send(msg)
    }
}
