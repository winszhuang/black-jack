package ws

import (
	"black-jack/card"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type NewClientRes struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (r *Room) OnJoin(c IClient) {
	r.mu.Lock()
	r.clients.Set(c, true)
	r.mu.Unlock()

	BroadcastGameSuccessRes(r, BroadcastJoin, NewClientRes{ID: c.GetID(), Name: c.GetLoginInfo().UserName}, fmt.Sprintf("玩家%s進入頻道", c.GetID()))
	BroadcastGameSuccessRes(r, UpdatePlayersDetail, r.getAllClientDetail(), "更新所有玩家資訊")
}

func (r *Room) OnLeave(c IClient) {
	r.mu.Lock()
	if _, ok := r.clients.Get(c); ok {
		c.CloseWsSend()
		r.clients.Delete(c)
	}
	r.mu.Unlock()

	// 某玩家斷線後若不到2人 遊戲中斷
	if r.clients.Len() < 2 && (r.isAllPlayerSameState(Play) || r.isAllPlayerSameState(Stop)) {
		r.Restart()
		BroadcastGameSuccessRes(r, BroadcastReStart, nil, fmt.Sprintf("玩家人數不到2人 遊戲重新"))
		BroadcastGameSuccessRes(r, UpdatePlayersDetail, r.getAllClientDetail(), "更新所有玩家資訊")
	} else {
		BroadcastGameSuccessRes(r, BroadcastLeave, c.GetID(), fmt.Sprintf("ClientID-%s玩家離開遊戲", c.GetID()))
		BroadcastGameSuccessRes(r, UpdatePlayersDetail, r.getAllClientDetail(), "更新所有玩家資訊")
	}
}

func (r *Room) OnReady(c IClient) {
	r.mu.Lock()
	notWaiting := c.GetCurrentState() > Wait
	if notWaiting {
		SendGameErrRes(c, ClientReady, ErrForWrongFlow, "錯誤的流程")
		r.mu.Unlock()
		return
	}

	c.UpdateCurrentState(Ready)
	r.mu.Unlock()

	SendGameSuccessRes(c, ClientReady, c.GetID(), fmt.Sprintf("你已準備"))
	BroadcastGameSuccessRes(r, BroadcastReady, c.GetID(), fmt.Sprintf("ClientID-%s玩家已經按下準備", c.GetID()))
	BroadcastGameSuccessRes(r, UpdatePlayersDetail, r.getAllClientDetail(), "更新所有玩家資料")

	// 需要大於一個玩家才能開始遊戲
	if r.isMoreThanOnePlayer() {
		r.checkAllPlayerReadyToStart()
	}
}

func (r *Room) onGameStart() {
	r.updateAllPlayerState(Play)
	err := r.buildAllPlayerCards()
	if err != nil {
		panic(err)
	}

	BroadcastGameSuccessRes(r, BroadcastGameStart, nil, "遊戲開始!!")
	BroadcastGameSuccessRes(r, UpdatePlayersDetail, r.getAllClientDetail(), "更新所有玩家資料")
}

type NewCardInfo struct {
	ClientID uuid.UUID `json:"client_id"`
	CardInfo card.Card `json:"card_info"`
}

func (r *Room) OnHit(c IClient) {
	r.mu.Lock()
	notPlaying := c.GetCurrentState() != Play
	if notPlaying {
		SendGameErrRes(c, ClientHit, ErrForWrongFlow, "錯誤的流程")
		r.mu.Unlock()
		return
	}

	// 發牌
	card, err := r.cardDealer.DealCard()
	if err != nil {
		SendGameErrRes(c, ClientHit, ErrForServerError, "伺服器問題 - 發牌錯誤")
		r.mu.Unlock()
		panic(err)
		return
	}

	// 更新牌給該玩家
	c.AddCard(card)
	r.mu.Unlock()

	result := NewCardInfo{
		ClientID: c.GetID(),
		CardInfo: card,
	}

	SendGameSuccessRes(c, ClientHit, result, fmt.Sprintf("你獲得新的一副牌"))
	BroadcastGameSuccessRes(r, BroadcastHit, result, fmt.Sprintf("ClientID-%s玩家獲得新牌", c.GetID()))
	BroadcastGameSuccessRes(r, UpdatePlayersDetail, r.getAllClientDetail(), "更新所有玩家資料")

	r.checkPlayerBustThenStop(c)
	r.checkAllPlayerStopThenEnd()
}

func (r *Room) OnStand(c IClient) {
	r.mu.Lock()
	notPlaying := c.GetCurrentState() != Play
	if notPlaying {
		r.mu.Unlock()
		SendGameErrRes(c, ClientStand, ErrForWrongFlow, "錯誤的流程")
		return
	}

	// 更新玩家狀態
	c.UpdateCurrentState(Stop)
	r.mu.Unlock()

	BroadcastGameSuccessRes(r, BroadcastStand, c.GetID(), fmt.Sprintf("ClientID-%s玩家停止要牌", c.GetID()))
	BroadcastGameSuccessRes(r, UpdatePlayersDetail, r.getAllClientDetail(), "更新所有玩家資料")

	r.checkAllPlayerStopThenEnd()
}

func (r *Room) onGameEnd() {
	r.updateAllPlayerState(End)

	winners, isExist := r.calculateFinalWinners()
	if !isExist {
		BroadcastGameSuccessRes(r, BroadcastGameOver, []uuid.UUID{}, "沒有任何玩家獲勝")
		return
	}

	winnerIds := []uuid.UUID{}
	message := ""
	for _, winner := range winners {
		id := winner.GetID()
		winnerIds = append(winnerIds, id)
		message += id.String() + " "
	}
	BroadcastGameSuccessRes(r, BroadcastGameOver, winnerIds, fmt.Sprintf("獲得勝利的是: ")+message)
	BroadcastGameSuccessRes(r, UpdatePlayersDetail, r.getAllClientDetail(), "更新所有玩家資料")

	time.Sleep(time.Second * 3)
	r.Restart()
	BroadcastGameSuccessRes(r, UpdatePlayersDetail, r.getAllClientDetail(), "更新所有玩家資料")
}
