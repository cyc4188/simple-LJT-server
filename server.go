package main

import (
	"LJT-server/proto"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

type GameServer struct {
	proto.UnimplementedGameServer
	game    *Game // TODO: support for multiple games match
	clients map[uuid.UUID]*Client
	mu      sync.RWMutex
}

func NewGameServer(game *Game) *GameServer {
	server := &GameServer{
		game:    game,
		clients: make(map[uuid.UUID]*Client),
	}
	go server.watchChange()
	return server
}

// 处理连接请求
// refer: proto ConnectRequest
func (s *GameServer) Connecting(ctx context.Context, req *proto.ConnectRequest) (*proto.ConnectResponse, error) {
	if s.game.isGameFull() {
		return nil, errors.New("server is full")
	}
	clientId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	// check if player already connected
	s.mu.RLock()
	if _, ok := s.clients[clientId]; ok {
		s.mu.RUnlock()
		return nil, errors.New("player already connected")
	}
	s.mu.RUnlock()

	newClient := NewClient(clientId, req.Name, s.game)

	log.Printf("player %s connected", req.Name)
	// add client to the server
	s.mu.Lock()
	s.clients[clientId] = newClient
	s.mu.Unlock()

	// add client to the game
	s.game.Mu.Lock()
	s.game.addClient(newClient)
	s.game.Mu.Unlock()
	if s.game.checkPlayerCount() {
		go s.game.startGame()
	}

	return &proto.ConnectResponse{
		Token:   "test",
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
	log.Printf("receive: %v", req)
	clientId, err := uuid.Parse(req.GetPlayCards().Player.GetId())
	// 根据请求的uuid得到对应的streamServer
	s.clients[clientId].streamServer = stream

	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				log.Printf("receive error: %v", err)
				return
			}
			log.Printf("receive: %v", req)
			// push to action chan
			switch req.GetRequest().(type) {
			case *proto.StreamRequest_PlayCards:
				s.handlePlayCardsRequest(req, clientId)
			case *proto.StreamRequest_Pass:
				s.handlePassRequest(req, clientId)
			}
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

func (s *GameServer) handlePlayCardsRequest(req *proto.StreamRequest, playerId uuid.UUID) {
	player, err := s.GetPlayerById(playerId)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	s.game.ActionChan <- PlayCards{
		player: player,
		cards:  CardsFromProto(req.GetPlayCards().Cards),
	}
}
func (s *GameServer) handlePassRequest(req *proto.StreamRequest, playerId uuid.UUID) {
	player, err := s.GetPlayerById(playerId)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	s.game.ActionChan <- Pass{
		player: player,
	}
}

func (s *GameServer) watchChange() {
	for {
		change := <-s.game.ChangeChan
		switch change.(type) {
		case GameStatus:
			change := change.(GameStatus)
			s.broadcast(change)
		case GameFail:
			change := change.(GameFail)
			s.sendFail(change)
		}
	}
}

// broadcast game status to all clients
// and its own cards
func (s *GameServer) broadcast(status GameStatus) {
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

		msg := status.ToProto(client.Player)

		// msg := &proto.StreamResponse{
		// 	Response: &proto.StreamResponse_Continue{
		// 		Continue: &proto.Continue{
		// 			Score:         int32(status.score),
		// 			Players:       players,
		// 			CurrentCards:  current_cards,
		// 			CurrentPlayer: status.current_player.ToProto(),
		// 			Cards:         cards,
		// 		},
		// 	},
		// }

		client.streamServer.Send(msg)
	}
}

// send play fail to client
func (s *GameServer) sendFail(fail GameFail) {
	for _, client := range s.clients {
		// check client and player
		if client != fail.player.client {
			continue
		}
		if client.streamServer == nil {
			log.Println("client ", client.uuid, " stream server is nil")
			continue
		}
		msg := fail.ToProto()
		client.streamServer.Send(msg)
	}
}

func (s *GameServer) GetPlayerById(id uuid.UUID) (*Player, error) {
	// id found in clients:
	if client, ok := s.clients[id]; ok {
		return client.Player, nil
	}
	return nil, fmt.Errorf("uuid: %v not found", id)
}
