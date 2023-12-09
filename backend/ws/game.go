package ws

import (
	"black-jack/game"
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"log"
	"sync"
	"time"
)

const MaxPoint = 21
const DealerClientID = "dealer"

type Game struct {
	clients *orderedmap.OrderedMap[IClient, bool] // 註冊的所有玩家

	cardDealer game.ICardDealer // 發牌員

	mu *sync.RWMutex // 鎖
}

func NewGame(cardDealer game.ICardDealer) *Game {
	g := &Game{
		clients:    orderedmap.New[IClient, bool](),
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
		client.InitPlayerInfo()
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

func (g *Game) checkPlayerBustThenStop(c IClient) {
	if g.isPlayerBust(c) {
		g.mu.Lock()
		c.UpdateCurrentState(Stop)
		g.mu.Unlock()

		BroadcastSuccessRes(g, BroadcastStand, c.GetID(), fmt.Sprintf("ClientID-%s玩家已經停止動作", c.GetID()))
		BroadcastSuccessRes(g, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")
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
func (g *Game) isPlayerBust(c IClient) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	totalPoints := c.CalculateTotalPoints()
	return totalPoints > MaxPoint
}

func (g *Game) updateAllPlayerState(state UserState) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		client.UpdateCurrentState(state)
	}
}

func (g *Game) calculateFinalWinners() ([]IClient, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// 玩家沒有1人以上就不用算了
	if g.clients.Len() < 2 {
		return []IClient{}, false
	}

	// 紀錄玩家得分-玩家資訊
	pointClientMap := map[int][]IClient{}
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		clientPoint := client.CalculateTotalPoints()
		// 爆點不理他
		if clientPoint > MaxPoint {
			continue
		}

		if _, isExist := pointClientMap[clientPoint]; !isExist {
			pointClientMap[clientPoint] = []IClient{client}
		} else {
			pointClientMap[clientPoint] = append(pointClientMap[clientPoint], client)
		}
	}

	// 算出最大點數
	maxPoint := 0
	for point, _ := range pointClientMap {
		if point > maxPoint {
			maxPoint = point
		}
	}
	if _, isExist := pointClientMap[maxPoint]; isExist {
		return pointClientMap[maxPoint], true
	}
	log.Println("靠北 這邊也會錯哦 算分", maxPoint)
	return []IClient{}, false
}

func (g *Game) Broadcast(data []byte) {
	g.mu.Lock()
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		client.WsSend(data)
	}
	g.mu.Unlock()
}

func (g *Game) getAllClientDetail() []ClientDetail {
	g.mu.Lock()
	defer g.mu.Unlock()

	var allClientsDetail []ClientDetail
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key

		allClientsDetail = append(allClientsDetail, client.GetGameDetail())
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
		if client.GetCurrentState() == state {
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
		client.AddCard(card1)
		client.AddCard(card2)
	}
	return nil
}

func (g *Game) getClient(clientID string) IClient {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		if client.GetID() == clientID {
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
