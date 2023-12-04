package ws

import (
	"black-jack/game"
	"fmt"
)

func (g *Game) OnRegister(c *Client) {
	g.mu.Lock()
	g.clients.Set(c, true)
	g.mu.Unlock()

	SendSuccessRes(c, SomeOneJoin, c.ID, fmt.Sprintf("你好 以下提供給你專屬ID"))
	BroadcastSuccessRes(c, SomeOneJoin, c.ID, fmt.Sprintf("玩家%s進入頻道", c.ID))
	BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資訊")
}

func (g *Game) OnUnRegister(c *Client) {
	g.mu.Lock()
	if _, ok := g.clients.Get(c); ok {
		g.clients.Delete(c)
		close(c.send)
	}
	g.mu.Unlock()

	BroadcastSuccessRes(c, SomeOneLeave, c.ID, fmt.Sprintf("ClientID-%s玩家離開遊戲", c.ID))
	BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資訊")
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

	BroadcastSuccessRes(c, SomeOneReady, c.ID, fmt.Sprintf("ClientID-%s玩家已經按下準備", c.ID))
	BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")

	g.checkAllReadyToStart()
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

type NewCardInfo struct {
	ClientID string    `json:"client_id"`
	CardInfo game.Card `json:"card_info"`
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
		ClientID: c.ID,
		CardInfo: card,
	}

	SendSuccessRes(c, SomeOneHit, result, fmt.Sprintf("你獲得新的一副牌"))
	BroadcastSuccessRes(c, SomeOneHit, result, fmt.Sprintf("ClientID-%s玩家獲得新牌", c.ID))
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
	BroadcastSuccessRes(c, SomeOneStand, c.ID, fmt.Sprintf("ClientID-%s玩家停止要牌", c.ID))
	BroadcastSuccessRes(c, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")

	g.checkAllStopToEnd()
}

func (g *Game) onGameEnd() {
	winner := g.calculateFinalWinner()
	g.updateAllPlayerState(End)

	g.Broadcast(WSResponse{
		MsgCode:   GameOver,
		Data:      winner,
		Success:   true,
		ErrorCode: 0,
		Message:   fmt.Sprintf("ClientID-%s玩家獲得勝利", winner.ID),
	}.Byte())
}
