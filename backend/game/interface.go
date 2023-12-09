package game

type ICardDealer interface {
	InitializeDeck()
	ShuffleDeck()
	DealCard() (Card, error)
}
