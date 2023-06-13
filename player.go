package main

import (
	"LJT-server/proto"
	"errors"
)

type Player struct {
    client *Client 
    Cards []Card
    score int
    index uint
}

func NewPlayer(client *Client) *Player {
    return &Player{
        client: client,
        Cards: make([]Card, 0),
        score: 0,
    }
}

func (player *Player) GetCardNum() int {
    return len(player.Cards)
}

// 是否有某张牌
func (player *Player) HasCard(card Card) bool {
    for _, c := range player.Cards {
        if c == card {
            return true
        }  
    }
    return false
}

// 是否有某些牌
func (player *Player) HasCards(cards []Card) bool {
    return IsSubsetForCards(player.Cards, cards) 
}

// 打出一些牌
func (player *Player) PlayCards(cards []Card) error {
    if !player.HasCards(cards) {
        return errors.New("player doesn't have these cards")
    }
    for _, card := range cards {
        for i, c := range player.Cards {
            if c == card {
                player.Cards = append(player.Cards[:i], player.Cards[i+1:]...)
                break
            }
        }
    }
    return nil
}

func (player *Player) ToProto() *proto.Player {
    return &proto.Player{
        Id: player.client.uuid.String(),
        Name: player.client.name,
        Score: int32(player.score),
        CardNum: int32(player.GetCardNum()),
        Index: uint32(player.index),
    }
}
