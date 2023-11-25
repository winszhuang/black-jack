package ws

import (
	"black-jack/game"
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const MaxPoint = 21
const DealerClientID = "dealer"

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Game maintains the set of active clients and broadcasts messages to the
// clients.
type Game struct {
	// Registered clients.
	clients *orderedmap.OrderedMap[*Client, bool]

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// ready requests from clients
	ready chan *Client

	// hit requests from clients
	hit chan *Client

	// stand requests from clients
	stand chan *Client

	start chan interface{}
	end   chan interface{}

	cardDealer *game.CardDealer
}

func NewGame(cardDealer *game.CardDealer) *Game {
	g := &Game{
		broadcast:  make(chan []byte),
		register:   make(chan *Client, 30),
		unregister: make(chan *Client, 30),
		ready:      make(chan *Client, 30),
		hit:        make(chan *Client, 30),
		stand:      make(chan *Client, 30),
		start:      make(chan interface{}),
		end:        make(chan interface{}),
		clients:    orderedmap.New[*Client, bool](),
		cardDealer: cardDealer,
	}
	// 創建一個莊家
	g.newDealerClient()
	g.Restart()
	return g
}

func (g *Game) Run() {
	go g.listenChanReceive()
	go g.checkGameFlow()
	go g.listenBroadcast()
}

func (g *Game) newDealerClient() {
	dealerClient := NewClient(g, nil, DealerClientID)
	g.clients.Set(dealerClient, true)
}

func (g *Game) listenChanReceive() {
	for {
		select {
		case client := <-g.register:
			g.onRegister(client)
		case client := <-g.unregister:
			g.onUnRegister(client)
		case client := <-g.ready:
			g.onReady(client)
		case client := <-g.hit:
			g.onHit(client)
		case client := <-g.stand:
			g.onStand(client)
		case <-g.start:
			g.onGameStart()
		case <-g.end:
			g.onGameEnd()
			//default:
		}
	}
}

func (g *Game) checkGameFlow() {
	for {
		if g.isAllPlayerReadyExceptDealer() {
			g.ready <- g.getClient(DealerClientID)
		}

		if g.isGameStart() {
			g.start <- true
		}

		if g.isGameEnd() {
			g.end <- true
		}

		if client, isExist := g.findSomeOneCrashPoint(); isExist {
			g.stand <- client
		}
	}
}

func (g *Game) updateAllPlayerState(state UserState) {
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		client.playInfo.currentState = state
	}
}

func (g *Game) listenBroadcast() {
	for {
		select {
		case message := <-g.broadcast:
			for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
				client := pair.Key
				select {
				case client.send <- message:
				default:
					close(client.send)
					g.clients.Delete(client)
				}
			}
		}
	}
}

func (g *Game) calculateFinalWinner() *Client {
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
	g.broadcast <- data
}

func (g *Game) NewClient(client *Client) {
	g.register <- client
}

func (g *Game) onRegister(c *Client) {
	g.clients.Set(c, true)
	SendSuccessRes(c, SomeOneJoin, c.ID, fmt.Sprintf("你好 以下提供給你專屬ID"))
	BroadcastSuccessRes(c, SomeOneJoin, g.getAllClientDetail(), fmt.Sprintf("ID-%s玩家進入遊戲", c.ID))
}

func (g *Game) onUnRegister(c *Client) {
	if _, ok := g.clients.Get(c); ok {
		g.clients.Delete(c)
		close(c.send)
	}

	BroadcastSuccessRes(c, SomeOneLeave, g.getAllClientDetail(), fmt.Sprintf("ID-%s玩家離開遊戲", c.ID))
}

type ClientDetail struct {
	ID   string    `json:"id"`
	Deck game.Deck `json:"deck"`
}

func (g *Game) getAllClientDetail() []ClientDetail {
	var allClientsDetail []ClientDetail
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key

		allClientsDetail = append(allClientsDetail, ClientDetail{
			ID:   client.ID,
			Deck: client.playInfo.deck,
		})
	}
	return allClientsDetail
}

type ReadyNotification struct {
	ID string `json:"id"`
}

func (g *Game) onReady(c *Client) {
	notWaiting := c.playInfo.currentState > Wait
	if notWaiting {
		SendErrRes(c, SomeOneReady, ErrForWrongFlow, "錯誤的流程")
		return
	}

	c.playInfo.currentState = Ready
	BroadcastSuccessRes(c, SomeOneReady, c.ID, fmt.Sprintf("ID-%s玩家已經按下準備", c.ID))
	BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")
}

func (g *Game) isGameStart() bool {
	count := 0
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		if client.playInfo.currentState == Ready {
			count++
		}
	}
	return count == g.clients.Len()
}

func (g *Game) isGameEnd() bool {
	count := 0
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		if client.playInfo.currentState == Stop {
			count++
		}
	}
	return count == g.clients.Len()
}

func (g *Game) findSomeOneCrashPoint() (*Client, bool) {
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
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

func (g *Game) onHit(c *Client) {
	notPlaying := c.playInfo.currentState != Play
	if notPlaying {
		SendErrRes(c, SomeOneHit, ErrForWrongFlow, "錯誤的流程")
		return
	}

	// 發牌
	card, err := g.cardDealer.DealCard()
	if err != nil {
		SendErrRes(c, SomeOneHit, ErrForServerError, "伺服器問題 - 發牌錯誤")
		panic(err)
		return
	}

	// 更新牌給該玩家
	c.playInfo.deck = c.playInfo.deck.AddCard(card)

	result := NewCardInfo{
		ID:       c.ID,
		CardInfo: card,
	}

	// 廣撥給所有玩家
	BroadcastSuccessRes(c, SomeOneHit, result, fmt.Sprintf("ID-%s玩家獲得新牌", c.ID))
	BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")
}

// Restart 重新開始遊戲
func (g *Game) Restart() {
	g.cardDealer.InitializeDeck()
	g.cardDealer.ShuffleDeck()
}

func (g *Game) onStand(c *Client) {
	notPlaying := c.playInfo.currentState != Play
	if notPlaying {
		SendErrRes(c, SomeOneHit, ErrForWrongFlow, "錯誤的流程")
		return
	}

	// 更新玩家狀態
	c.playInfo.currentState = Stop

	// 廣撥給所有玩家
	BroadcastSuccessRes(c, SomeOneStand, c.ID, fmt.Sprintf("ID-%s玩家停止要牌", c.ID))
	BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")
}

func (g *Game) buildAllPlayerCards() error {
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

func (g *Game) RequestReady(c *Client) {
	g.ready <- c
}

func (g *Game) RequestHit(c *Client) {
	g.hit <- c
}

func (g *Game) RequestStand(c *Client) {
	g.stand <- c
}

func (g *Game) isAllPlayerReadyExceptDealer() bool {
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
		if client.playInfo.currentState == Ready {
			list = append(list, true)
		}
	}
	return len(list) == playerCount
}

func (g *Game) getClient(clientID string) *Client {
	for pair := g.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		if client.ID == clientID {
			return client
		}
	}
	return nil
}
