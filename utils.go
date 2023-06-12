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
