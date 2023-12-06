package game

type Card struct {
	Name   string `json:"name"`
	Suit   Suit   `json:"suit"`
	Symbol string `json:"symbol"`
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
