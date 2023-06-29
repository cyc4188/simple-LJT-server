package main

import (
	"math/rand"
	"time"
)

type GameRule interface {
	PlayerCount() int // 玩家数量
	DeckCount() int   // 几副牌
	CardCount() int
	CardPerPlayer() int              // 每个玩家的牌数
	PokerHands([]Card) bool          // 可打出的牌型
	CompareHands([]Card, []Card) int // 比较两副牌的大小
	generateDeck() []Card            // 生成牌组
	checkHandsIsValid([]Card) bool   // 检查牌组是否合法
}

type DdzGameRule struct {
}

const (
	DDZ_PLAYER_COUNT = 2 // test
	DDZ_DECK_COUNT   = 1
)

func NewDdzGameRule() *DdzGameRule {
	return &DdzGameRule{}
}

func (rule *DdzGameRule) PlayerCount() int {
	return DDZ_PLAYER_COUNT
}

func (rule *DdzGameRule) DeckCount() int {
	return DDZ_DECK_COUNT
}

func (rule *DdzGameRule) CardPerPlayer() int {
	return rule.CardCount() / rule.PlayerCount()
}

func (rule *DdzGameRule) CardCount() int {
	return rule.DeckCount() * CARD_PER_DECK
}

func (rule *DdzGameRule) PokerHands(cards []Card) bool {
	// TODO: check if cards is valid
	return true
}

// return 1 if cards1 > cards2
// return 0 if cards1 == cards2
// return -1 if cards1 < cards2
func (rule *DdzGameRule) CompareHands(cards1 []Card, cards2 []Card) int {
	// TODO: get order
	if len(cards2) == 0 {
		return 1
	}
	return cards1[0].Compare(&cards2[0])
}

// randomly generate decks and shuffle
func (rule *DdzGameRule) generateDeck() []Card {
	deckCount := rule.DeckCount()
	deck := make([]Card, 0, deckCount*CARD_PER_DECK)
	for i := 0; i < deckCount; i++ {
		deck = append(deck, generateDeck()...)
	}

	// shuffle the deck
	rand.Seed(time.Now().UnixNano())
	for i := range deck {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	return deck
}

// check the cards is valid
// TODO
func (rule *DdzGameRule) checkHandsIsValid(cards []Card) bool {
	return true
}
