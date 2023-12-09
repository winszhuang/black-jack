package card

import (
	"strconv"
)

type Card struct {
	Name   string `json:"name"` // e.g. "AC", "2D", "KS", "5H"
	Suit   Suit   `json:"suit"`
	Symbol string `json:"symbol"` // 數字的符號 e.g. 13是"K" 12是"Q" 4就是"4"
	Value  int    `json:"value"`
}

// Suit represents a card suit.
type Suit string

const (
	Hearts   Suit = "H"
	Diamonds      = "D"
	Clubs         = "C"
	Spades        = "S"
)

func NewCardByName(cardStr string) Card {
	if len(cardStr) < 2 {
		panic("invalid card string")
	}

	symbol := cardStr[:len(cardStr)-1]
	suit := Suit(cardStr[len(cardStr)-1:])

	var value int
	var err error
	if symbol == "A" {
		value = 11 // "A" 可以等於1或11
	} else if symbol == "K" || symbol == "Q" || symbol == "J" {
		value = 10 // "K", "Q", "J" 都等於10
	} else {
		value, err = strconv.Atoi(symbol)
		if err != nil {
			panic("invalid card value")
		}
	}

	return Card{
		Name:   cardStr,
		Suit:   suit,
		Symbol: symbol,
		Value:  value,
	}
}
