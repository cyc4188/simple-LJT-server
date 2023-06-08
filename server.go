package main

import (
	"LJT-server/proto"
	"context"
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
)

type GameServer struct {
    proto.UnimplementedGameServer;
    game *Game; // TODO: support for multiple games match
    clients map[uuid.UUID]*Client;
    mu sync.RWMutex;
}

func NewGameServer(game *Game) *GameServer {
    return &GameServer{
        game: game,
        clients: make(map[uuid.UUID]*Client),
    }
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
        s.broadcastGameStatus()
    }

    return &proto.ConnectResponse{
        Token: "test",
        Players: []*proto.Player{},
    }, nil
}

func (s *GameServer) Stream(stream proto.Game_StreamServer) error {
    ctx := stream.Context()

    go func() {
        // TODO: listen to client's stream
        for {
            req, err := stream.Recv()
            if err != nil {
                log.Printf("receive error: %v", err)
                return 
            }
            log.Printf("receive: %v", req.PlayCards) 
            s.broadcastGameStatus()
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

func (s *GameServer) broadcast(msg *proto.StreamResponse) {
    for _, client := range s.clients {
        if client.streamServer == nil {
            continue
        }
        client.streamServer.Send(msg)
    }
}

func (s *GameServer) broadcastGameStatus() { // 游戏进行中，广播游戏状态
    // TODO: broadcast game status
    for {
        select {
        case change := <-s.game.ChangeChan:
            // TODO
        default:

        }
    }
}
