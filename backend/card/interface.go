package card

type ICardDealer interface {
	InitializeDeck()
	ShuffleDeck()
	DealCard() (Card, error)
}
