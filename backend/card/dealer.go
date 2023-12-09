package card

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

var (
	suits  = []Suit{Hearts, Diamonds, Clubs, Spades}
	values = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
)

// CardDealer represents a card dealer.
type CardDealer struct {
	rand *rand.Rand
	Deck Deck // 牌組
}

// NewCardDealer creates a new card dealer with a shuffled deck.
func NewCardDealer(rand *rand.Rand) *CardDealer {
	return &CardDealer{rand: rand}
}

// InitializeDeck initializes the deck of cards.
func (cd *CardDealer) InitializeDeck() {
	cd.Deck = Deck{}
	for _, suit := range suits {
		for _, value := range values {
			var cardValue int
			if value == "A" {
				cardValue = 11 // "A" 可以等於1或11
			} else if value == "K" || value == "Q" || value == "J" {
				cardValue = 10 // "K", "Q", "J" 都等於10
			} else {
				cardValue, _ = strconv.Atoi(value)
			}

			card := Card{
				Name:   value + string(suit),
				Suit:   suit,
				Symbol: value,
				Value:  cardValue,
			}
			cd.Deck = append(cd.Deck, card)
		}
	}
}

// ShuffleDeck shuffles the deck of cards.
func (cd *CardDealer) ShuffleDeck() {
	cd.rand.Seed(time.Now().UnixNano())
	cd.rand.Shuffle(len(cd.Deck), func(i, j int) {
		cd.Deck[i], cd.Deck[j] = cd.Deck[j], cd.Deck[i]
	})
}

// DealCard deals a card from the deck.
func (cd *CardDealer) DealCard() (Card, error) {
	if len(cd.Deck) == 0 {
		return Card{}, errors.New("No cards left in the deck.")
	}
	card := cd.Deck[0]
	cd.Deck = cd.Deck[1:]
	return card, nil
}
