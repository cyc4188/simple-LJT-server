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
    game *Game;
    clients map[uuid.UUID]*Client;
    mu sync.RWMutex;
}

func NewGameServer() *GameServer {
    return &GameServer{
        game: NewGame(NewDdzGameRule()),
        clients: make(map[uuid.UUID]*Client),
    }
}

// 处理连接请求
// refer: proto ConnectRequest
func (s *GameServer) Connect(ctx context.Context, req *proto.ConnectRequest) (*proto.ConnectResponse, error) {
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

    return nil, nil
}

func (s *GameServer) Stream(stream proto.Game_StreamServer) error {
    // TODO: not implemented
    return nil
}
