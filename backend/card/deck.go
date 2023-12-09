package card

type Deck []Card

func NewDeck() []Card {
	return []Card{}
}

func (d Deck) AddCard(card Card) Deck {
	d = append(d, card)
	return d
}

func (d Deck) CalculateTotalPoints() int {
	total := 0
	for _, card := range d {
		total += card.Value
	}
	return total
}
