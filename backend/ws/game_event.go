package ws

import (
	"black-jack/game"
	"fmt"
	"time"
)

func (g *Game) OnJoin(c IClient) {
	g.mu.Lock()
	g.clients.Set(c, true)
	g.mu.Unlock()

	SendSuccessRes(c, ClientJoin, c.GetID(), fmt.Sprintf("你好 以下提供給你專屬ID"))
	BroadcastSuccessRes(g, BroadcastJoin, c.GetID(), fmt.Sprintf("玩家%s進入頻道", c.GetID()))
	BroadcastSuccessRes(g, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資訊")
}

func (g *Game) OnLeave(c IClient) {
	g.mu.Lock()
	if _, ok := g.clients.Get(c); ok {
		c.CloseWsSend()
		g.clients.Delete(c)
	}
	g.mu.Unlock()

	// 某玩家斷線後若不到2人 遊戲中斷
	if g.clients.Len() < 2 && (g.isAllPlayerSameState(Play) || g.isAllPlayerSameState(Stop)) {
		g.Restart()
		BroadcastSuccessRes(g, BroadcastReStart, nil, fmt.Sprintf("玩家人數不到2人 遊戲重新"))
		BroadcastSuccessRes(g, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資訊")
	} else {
		BroadcastSuccessRes(g, BroadcastLeave, c.GetID(), fmt.Sprintf("ClientID-%s玩家離開遊戲", c.GetID()))
		BroadcastSuccessRes(g, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資訊")
	}
}

func (g *Game) OnReady(c IClient) {
	g.mu.Lock()
	notWaiting := c.GetCurrentState() > Wait
	if notWaiting {
		SendErrRes(c, ClientReady, ErrForWrongFlow, "錯誤的流程")
		g.mu.Unlock()
		return
	}

	c.UpdateCurrentState(Ready)
	g.mu.Unlock()

	SendSuccessRes(c, ClientReady, c.GetID(), fmt.Sprintf("你已準備"))
	BroadcastSuccessRes(g, BroadcastReady, c.GetID(), fmt.Sprintf("ClientID-%s玩家已經按下準備", c.GetID()))
	BroadcastSuccessRes(g, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")

	// 需要大於一個玩家才能開始遊戲
	if g.isMoreThanOnePlayer() {
		g.checkAllPlayerReadyToStart()
	}
}

func (g *Game) onGameStart() {
	g.updateAllPlayerState(Play)
	err := g.buildAllPlayerCards()
	if err != nil {
		panic(err)
	}

	g.Broadcast(WSResponse{
		MsgCode:   BroadcastGameStart,
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

func (g *Game) OnHit(c IClient) {
	g.mu.Lock()
	notPlaying := c.GetCurrentState() != Play
	if notPlaying {
		SendErrRes(c, ClientHit, ErrForWrongFlow, "錯誤的流程")
		g.mu.Unlock()
		return
	}

	// 發牌
	card, err := g.cardDealer.DealCard()
	if err != nil {
		SendErrRes(c, ClientHit, ErrForServerError, "伺服器問題 - 發牌錯誤")
		g.mu.Unlock()
		panic(err)
		return
	}

	// 更新牌給該玩家
	c.AddCard(card)
	g.mu.Unlock()

	result := NewCardInfo{
		ClientID: c.GetID(),
		CardInfo: card,
	}

	SendSuccessRes(c, ClientHit, result, fmt.Sprintf("你獲得新的一副牌"))
	BroadcastSuccessRes(g, BroadcastHit, result, fmt.Sprintf("ClientID-%s玩家獲得新牌", c.GetID()))
	BroadcastSuccessRes(g, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")

	g.checkPlayerBustThenStop(c)
	g.checkAllPlayerStopThenEnd()
}

func (g *Game) OnStand(c IClient) {
	g.mu.Lock()
	notPlaying := c.GetCurrentState() != Play
	if notPlaying {
		g.mu.Unlock()
		SendErrRes(c, ClientStand, ErrForWrongFlow, "錯誤的流程")
		return
	}

	// 更新玩家狀態
	c.UpdateCurrentState(Stop)
	g.mu.Unlock()

	BroadcastSuccessRes(g, BroadcastStand, c.GetID(), fmt.Sprintf("ClientID-%s玩家停止要牌", c.GetID()))
	BroadcastSuccessRes(g, UpdatePlayersDetail, g.getAllClientDetail(), "更新所有玩家資料")

	g.checkAllPlayerStopThenEnd()
}

func (g *Game) onGameEnd() {
	g.updateAllPlayerState(End)

	winners, isExist := g.calculateFinalWinners()
	if !isExist {
		g.Broadcast(WSResponse{
			MsgCode:   BroadcastGameOver,
			Data:      nil,
			Success:   true,
			ErrorCode: 0,
			Message:   fmt.Sprintf("沒有任何玩家獲勝"),
		}.Byte())
		return
	}

	winnerIds := []string{}
	message := ""
	for _, winner := range winners {
		id := winner.GetID()
		winnerIds = append(winnerIds, id)
		message += id + " "
	}
	g.Broadcast(WSResponse{
		MsgCode:   BroadcastGameOver,
		Data:      winnerIds,
		Success:   true,
		ErrorCode: 0,
		Message:   fmt.Sprintf("獲得勝利的是: ") + message,
	}.Byte())

	g.Broadcast(WSResponse{
		MsgCode:   UpdatePlayersDetail,
		Data:      g.getAllClientDetail(),
		Success:   true,
		ErrorCode: 0,
		Message:   "更新所有玩家資訊",
	}.Byte())

	time.Sleep(time.Second * 3)
	g.Restart()
	g.Broadcast(WSResponse{
		MsgCode:   UpdatePlayersDetail,
		Data:      g.getAllClientDetail(),
		Success:   true,
		ErrorCode: 0,
		Message:   "更新所有玩家資訊",
	}.Byte())
}
