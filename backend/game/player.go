package game

type Player struct {
	Deck Deck
	ID   string
}

func NewPlayer(playerID string) *Player {
	return &Player{Deck: NewDeck(), ID: playerID}
}

func (p *Player) AddCard(card Card) {
	cards := p.Deck.AddCard(card)
	p.Deck = cards
}
