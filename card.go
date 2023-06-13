package main

import "LJT-server/proto"

const (
    CARD_PER_DECK = 54
)

// poker card
type Card struct {
    suit int // 0-3
    rank int // 1-14
}

var suits = []string{"spade", "heart", "diamond", "club"}
var mapRank = map[int]string{
    1: "3",
    2: "4",
    3: "5",
    4: "6",
    5: "7",
    6: "8",
    7: "9",
    8: "10",
    9: "J",
    10: "Q",
    11: "K",
    12: "A",
    13: "2",
    14: "Joker0",
    15: "Joker1",
}

func (c *Card) Compare(other *Card) int {
    if c.rank > other.rank {
        return 1
    }
    if c.rank < other.rank {
        return -1
    }
    return 0
}

func (c *Card) String() string {
    return suits[c.suit] + "." + mapRank[c.rank]
}


func (c *Card) ToProto() *proto.Card {
    return &proto.Card {
        Suit: int32(c.suit), 
        Rank: int32(c.rank),
    }
}

// generate a deck of cards
func generateDeck() []Card {
    deck := make([]Card, 0, CARD_PER_DECK) 
    for i := 0; i < 4; i++ {
        for j := 1; j <= 13; j++ {
            deck = append(deck, Card{i, j})
        }
    }
    deck = append(deck, Card{4, 14})
    deck = append(deck, Card{4, 15})
    return deck
}

func IsSubsetForCards(cards []Card, subset []Card) bool {
    set := make(map[Card]int)
    for _, card := range cards {
        set[card]++
    }
    for _, card := range subset {
        if set[card] == 0 {
            return false
        }
        set[card]--
    }
    return true
}
