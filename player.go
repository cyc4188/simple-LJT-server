package main

import "LJT-server/proto"

type Player struct {
    client *Client 
    Cards []Card
    score int
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

func (player *Player) ToProto() *proto.Player {
    return &proto.Player{
        Id: player.client.uuid.String(),
        Name: player.client.name,
        Score: int32(player.score),
        CardNum: int32(player.GetCardNum()),
    }
}
