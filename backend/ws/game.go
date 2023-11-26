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
	// 創建一個莊家
	g.newDealerClient()
	g.Restart()
	return g
}

func (g *Game) newDealerClient() {
	dealerClient := NewClient(g, nil, DealerClientID)
	g.clients.Set(dealerClient, true)
}

// Restart 重新開始遊戲
func (g *Game) Restart() {
	g.cardDealer.InitializeDeck()
	g.cardDealer.ShuffleDeck()
}

func (g *Game) checkAllReadyToStart() {
	if !g.isAllPlayerSameStateExceptDealer(Ready) {
		return
	}

	go func() {
		// 等待2秒，莊家變成準備模式
		time.Sleep(time.Second * 2)

		// dealer更新狀態為ready
		dealer := g.getClient(DealerClientID)
		g.mu.Lock()
		dealer.playInfo.currentState = Ready
		g.mu.Unlock()

		BroadcastSuccessRes(dealer, SomeOneReady, dealer.ID, fmt.Sprintf("ID-%s莊家已經按下準備", dealer.ID))
		BroadcastSuccessRes(dealer, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")

		// 開始遊戲
		g.onGameStart()
	}()
}

func (g *Game) checkPlayerCrashPointThenStop(c *Client) {
	if g.isPlayerCrashPoint(c) {
		g.mu.Lock()
		c.playInfo.currentState = Stop
		g.mu.Unlock()

		BroadcastSuccessRes(c, SomeOneStand, c.ID, fmt.Sprintf("ID-%s玩家已經停止動作", c.ID))
		BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")
	}
}

func (g *Game) isPlayerCrashPoint(c *Client) bool {
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

func (g *Game) OnRegister(c *Client) {
	g.mu.Lock()
	g.clients.Set(c, true)
	g.mu.Unlock()

	SendSuccessRes(c, SomeOneJoin, c.ID, fmt.Sprintf("你好 以下提供給你專屬ID"))
	BroadcastSuccessRes(c, SomeOneJoin, g.getAllClientDetail(), fmt.Sprintf("ID-%s玩家進入遊戲", c.ID))
}

func (g *Game) OnUnRegister(c *Client) {
	g.mu.Lock()
	if _, ok := g.clients.Get(c); ok {
		g.clients.Delete(c)
		close(c.send)
	}
	g.mu.Unlock()

	BroadcastSuccessRes(c, SomeOneLeave, g.getAllClientDetail(), fmt.Sprintf("ID-%s玩家離開遊戲", c.ID))
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

func (g *Game) OnReady(c *Client) {
	g.mu.Lock()
	notWaiting := c.playInfo.currentState > Wait
	if notWaiting {
		SendErrRes(c, SomeOneReady, ErrForWrongFlow, "錯誤的流程")
		g.mu.Unlock()
		return
	}

	c.playInfo.currentState = Ready
	g.mu.Unlock()

	BroadcastSuccessRes(c, SomeOneReady, c.ID, fmt.Sprintf("ID-%s玩家已經按下準備", c.ID))
	BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")

	g.checkAllReadyToStart()
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

func (g *Game) findSomeOneCrashPoint() (*Client, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		if client.playInfo.currentState == Stop {
			continue
		}

		totalPoints := client.playInfo.deck.CalculateTotalPoints()
		if totalPoints > MaxPoint {
			return client, true
		}
	}
	return nil, false
}

type NewCardInfo struct {
	ID       string    `json:"id"`
	CardInfo game.Card `json:"cardInfo"`
}

func (g *Game) OnHit(c *Client) {
	g.mu.Lock()
	notPlaying := c.playInfo.currentState != Play
	if notPlaying {
		SendErrRes(c, SomeOneHit, ErrForWrongFlow, "錯誤的流程")
		g.mu.Unlock()
		return
	}

	// 發牌
	card, err := g.cardDealer.DealCard()
	if err != nil {
		SendErrRes(c, SomeOneHit, ErrForServerError, "伺服器問題 - 發牌錯誤")
		g.mu.Unlock()
		panic(err)
		return
	}

	// 更新牌給該玩家
	c.playInfo.deck = c.playInfo.deck.AddCard(card)
	g.mu.Unlock()

	result := NewCardInfo{
		ID:       c.ID,
		CardInfo: card,
	}

	// 廣撥給所有玩家
	BroadcastSuccessRes(c, SomeOneHit, result, fmt.Sprintf("ID-%s玩家獲得新牌", c.ID))
	BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")

	g.checkPlayerCrashPointThenStop(c)
	g.checkAllStopToEnd()
}

func (g *Game) OnStand(c *Client) {
	g.mu.Lock()
	notPlaying := c.playInfo.currentState != Play
	if notPlaying {
		g.mu.Unlock()
		SendErrRes(c, SomeOneStand, ErrForWrongFlow, "錯誤的流程")
		return
	}

	// 更新玩家狀態
	c.playInfo.currentState = Stop
	g.mu.Unlock()

	// 廣撥給所有玩家
	BroadcastSuccessRes(c, SomeOneStand, c.ID, fmt.Sprintf("ID-%s玩家停止要牌", c.ID))
	BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")

	g.checkAllStopToEnd()
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

func (g *Game) onGameStart() {
	g.updateAllPlayerState(Play)
	err := g.buildAllPlayerCards()
	if err != nil {
		panic(err)
	}

	g.Broadcast(WSResponse{
		MsgCode:   GameStart,
		Data:      true,
		Success:   true,
		ErrorCode: 0,
		Message:   "遊戲開始!!",
	}.Byte())
	g.Broadcast(WSResponse{
		MsgCode:   UpdatePlayersDetail,
		Data:      g.getAllClientDetail(),
		Success:   true,
		ErrorCode: 0,
		Message:   "更新所有玩家資料",
	}.Byte())
}

func (g *Game) onGameEnd() {
	winner := g.calculateFinalWinner()
	g.updateAllPlayerState(End)

	g.Broadcast(WSResponse{
		MsgCode:   GameOver,
		Data:      winner,
		Success:   true,
		ErrorCode: 0,
		Message:   fmt.Sprintf("ID-%s玩家獲得勝利", winner.ID),
	}.Byte())
}

func (g *Game) isAllPlayerSameStateExceptDealer(state UserState) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	playerCount := g.clients.Len() - 1
	if playerCount == 0 {
		return false
	}

	var list []bool
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		if client.ID == DealerClientID {
			continue
		}
		if client.playInfo.currentState == state {
			list = append(list, true)
		}
	}
	return len(list) == playerCount
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

func (g *Game) checkAllStopToEnd() {
	if !g.isAllPlayerSameStateExceptDealer(Stop) {
		return
	}

	go func() {
		// 等待2秒，莊家變成準備模式
		time.Sleep(time.Second * 1)

		// dealer更新狀態為stop
		dealer := g.getClient(DealerClientID)
		g.mu.Lock()
		dealer.playInfo.currentState = Stop
		g.mu.Unlock()

		BroadcastSuccessRes(dealer, SomeOneStand, dealer.ID, fmt.Sprintf("ID-%s莊家已經停止動作", dealer.ID))
		BroadcastSuccessRes(dealer, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")

		// 結束遊戲
		g.onGameEnd()
	}()
}
