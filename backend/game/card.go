package game

type Card struct {
	Suit   Suit   `json:"suit"`
	Symbol string `json:"symbol"`
	Value  int    `json:"value"`
}

// Suit represents a card suit.
type Suit int

const (
	Hearts Suit = iota + 1
	Diamonds
	Clubs
	Spades
)
