package game

type CardDealerMock struct {
	Deck Deck
}

func NewCardDealerMock(deck Deck) ICardDealer {
	return &CardDealerMock{Deck: deck}
}

func (c *CardDealerMock) InitializeDeck() {
	// nothing
}

func (c *CardDealerMock) ShuffleDeck() {
	// nothing
}

func (c *CardDealerMock) DealCard() (Card, error) {
	card := c.Deck[0]
	c.Deck = c.Deck[1:]
	return card, nil
}
