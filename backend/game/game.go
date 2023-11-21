package game

import (
	"errors"
	"sync"
)

type Game struct {
	mu          sync.Mutex // 用於鎖定並發操作的互斥鎖
	dealer      *CardDealer
	players     map[string]*Player
	isGameStart bool
}

func NewGame(dealer *CardDealer) *Game {
	game := &Game{
		dealer:      dealer,
		players:     make(map[string]*Player),
		mu:          sync.Mutex{},
		isGameStart: false,
	}

	game.Init()
	return game
}

func (g *Game) IsGameStart() bool {
	return g.isGameStart
}

func (g *Game) Init() {
	g.dealer.InitializeDeck()
	g.dealer.ShuffleDeck()
}

func (g *Game) AddPlayer(playerID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, isExist := g.players[playerID]; !isExist {
		g.players[playerID] = NewPlayer(playerID)
	}
}

func (g *Game) BuildPlayersCards() error {
	for _, player := range g.players {
		card1, err := g.dealer.DealCard()
		card2, err2 := g.dealer.DealCard()
		if err != nil {
			return err
		}
		if err2 != nil {
			return err2
		}
		// 一人兩張
		player.AddCard(card1)
		player.AddCard(card2)
	}
	return nil
}

func (g *Game) GetPlayersCards() map[string]Deck {
	result := map[string]Deck{}
	for id, player := range g.players {
		result[id] = player.Deck
	}
	return result
}

func (g *Game) getPlayer(playerID string) (*Player, error) {
	if _, isExist := g.players[playerID]; !isExist {
		return nil, errors.New("無此玩家 不可取得牌組")
	}
	return g.players[playerID], nil
}

func (g *Game) DealCardToPlayer(playerID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	player, err := g.getPlayer(playerID)
	if err != nil {
		return err
	}

	card, err := g.dealer.DealCard()
	if err != nil {
		return err
	}

	player.AddCard(card)
	return nil
}
