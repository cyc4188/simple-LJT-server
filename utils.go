package main

import "LJT-server/proto"

func CardFromProto(card *proto.Card) Card {
	return Card{
		suit: int(card.Suit),
		rank: int(card.Rank),
	}
}

func CardsFromProto(cards []*proto.Card) []Card {
	result := make([]Card, 0)
	for _, card := range cards {
		result = append(result, CardFromProto(card))
	}
	return result
}

func CardToProto(card Card) *proto.Card {
	return &proto.Card{
		Suit: int32(card.suit),
		Rank: int32(card.rank),
	}
}

func CardsToProto(cards []Card) []*proto.Card {
	result := make([]*proto.Card, 0)
	for _, card := range cards {
		result = append(result, CardToProto(card))
	}
	return result
}

func (gamestatus *GameStatus) ToProto(player *Player) *proto.StreamResponse {
	players := make([]*proto.Player, 0)
	for _, player := range gamestatus.Players {
		players = append(players, player.ToProto())
	}
	current_cards := make([]*proto.Card, 0)
	for _, card := range gamestatus.current_cards {
		current_cards = append(current_cards, card.ToProto())
	}

	last_played := make([]*proto.LastPlayed, 0)
	for _, player := range gamestatus.Players {
		last_played = append(last_played, &proto.LastPlayed{
			Player: player.ToProto(),
			Cards:  CardsToProto(player.LastPlayed),
		})
	}
	return &proto.StreamResponse{
		Response: &proto.StreamResponse_Continue{
			Continue: &proto.Continue{
				Score:         int32(gamestatus.score),
				CurrentCards:  current_cards,
				CurrentPlayer: gamestatus.current_player.ToProto(),
				Players:       last_played,
				Cards:         CardsToProto(player.Cards),
			},
		},
	}
}

func (gameFail *GameFail) ToProto() *proto.StreamResponse {
	return &proto.StreamResponse{
		Response: &proto.StreamResponse_Fail{
			Fail: &proto.Fail{
				Reason: gameFail.msg,
			},
		},
	}
}
