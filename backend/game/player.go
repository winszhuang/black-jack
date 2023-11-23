package game

type Player struct {
	Deck            Deck
	ID              string
	isActionStopped bool
}

func NewPlayer(playerID string) *Player {
	return &Player{Deck: NewDeck(), ID: playerID}
}

func (p *Player) AddCard(card Card) {
	cards := p.Deck.AddCard(card)
	p.Deck = cards
}

// 結束當局
func (p *Player) EndTurn() {
	p.isActionStopped = true
}

// 玩家是否停止動作
func (p *Player) IsActionStopped() bool {
	return p.isActionStopped
}

func (p *Player) CalculateCardsPoint() int {
	return p.Deck.CalculateTotalPoints()
}
