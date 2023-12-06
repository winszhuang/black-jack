package ws

import (
	"black-jack/game"
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"sync"
	"time"
)

const MaxPoint = 21
const DealerClientID = "dealer"

type Game struct {
	// 註冊的所有玩家
	clients *orderedmap.OrderedMap[*Client, bool]

	// 發牌員
	cardDealer *game.CardDealer

	// 鎖
	mu *sync.RWMutex
}

func NewGame(cardDealer *game.CardDealer) *Game {
	g := &Game{
		clients:    orderedmap.New[*Client, bool](),
		cardDealer: cardDealer,
		mu:         &sync.RWMutex{},
	}
	g.Restart()
	return g
}

// Restart 重新開始遊戲
func (g *Game) Restart() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.cardDealer.InitializeDeck()
	g.cardDealer.ShuffleDeck()
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		client.playInfo.Init()
	}
}

func (g *Game) checkAllPlayerReadyToStart() {
	if !g.isAllPlayerSameState(Ready) {
		return
	}

	go func() {
		// 等待1秒，開始遊戲
		time.Sleep(time.Second)
		// 開始遊戲
		g.onGameStart()
	}()
}

func (g *Game) checkPlayerBustThenStop(c *Client) {
	if g.isPlayerBust(c) {
		g.mu.Lock()
		c.playInfo.currentState = Stop
		g.mu.Unlock()

		BroadcastSuccessRes(c, BroadcastStand, c.ID, fmt.Sprintf("ClientID-%s玩家已經停止動作", c.ID))
		BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")
	}
}

func (g *Game) checkAllPlayerStopThenEnd() {
	if !g.isAllPlayerSameState(Stop) {
		return
	}

	go func() {
		// 等待1秒，準備結束遊戲
		time.Sleep(time.Second * 1)
		// 結束遊戲
		g.onGameEnd()
	}()
}

// 玩家是否爆牌
func (g *Game) isPlayerBust(c *Client) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	totalPoints := c.playInfo.deck.CalculateTotalPoints()
	return totalPoints > MaxPoint
}

func (g *Game) updateAllPlayerState(state UserState) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		client.playInfo.currentState = state
	}
}

func (g *Game) calculateFinalWinner() *Client {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// #TODO
	// 有可能同點數
	// 兩個人都爆點數
	// 多人都爆點數
	firstClient := g.clients.Oldest().Key
	currMaxPoint := firstClient.playInfo.deck.CalculateTotalPoints()
	winnerClient := firstClient
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		if client == firstClient {
			continue
		}

		clientPoint := client.playInfo.deck.CalculateTotalPoints()
		if clientPoint > currMaxPoint {
			winnerClient = client
			currMaxPoint = clientPoint
		}
	}
	return winnerClient
}

func (g *Game) Broadcast(data []byte) {
	g.mu.Lock()
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		select {
		case client.send <- data:
		default:
			close(client.send)
			g.clients.Delete(client)
		}
	}
	g.mu.Unlock()
}

type ClientDetail struct {
	ID    string    `json:"id"`
	Deck  game.Deck `json:"deck"`
	State UserState `json:"state"`
}

func (g *Game) getAllClientDetail() []ClientDetail {
	g.mu.Lock()
	defer g.mu.Unlock()

	var allClientsDetail []ClientDetail
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key

		allClientsDetail = append(allClientsDetail, ClientDetail{
			ID:    client.ID,
			Deck:  client.playInfo.deck,
			State: client.playInfo.currentState,
		})
	}
	return allClientsDetail
}

type ReadyNotification struct {
	ID string `json:"id"`
}

func (g *Game) isAllPlayerSameState(state UserState) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	count := 0
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		if client.playInfo.currentState == state {
			count++
		}
	}
	return count == g.clients.Len()
}

func (g *Game) buildAllPlayerCards() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		card1, err := g.cardDealer.DealCard()
		card2, err2 := g.cardDealer.DealCard()
		if err != nil {
			return err
		}
		if err2 != nil {
			return err2
		}
		client.playInfo.deck = client.playInfo.deck.AddCard(card1)
		client.playInfo.deck = client.playInfo.deck.AddCard(card2)
	}
	return nil
}

func (g *Game) getClient(clientID string) *Client {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		if client.ID == clientID {
			return client
		}
	}
	return nil
}

func (g *Game) isMoreThanOnePlayer() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.clients.Len() > 1
}
