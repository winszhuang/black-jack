package ws

import (
	"black-jack/game"
	"fmt"
	"math/rand"
)

const (
	MaxPoint = 21
)

var GameEngine *game.Game

func init() {
	cardDealer := game.NewCardDealer(rand.New(rand.NewSource(848484814486486)))
	GameEngine = game.NewGame(cardDealer)
}

// msg_code 1
func OnStart(c *Client) {
	if GameEngine.IsGameStart() {
		SendErrRes(c, Start, ErrForNotCorrectState, "遊戲已開始 不得重新開始遊戲")
		return
	}

	GameEngine.Init()
	for client, isExist := range c.Hub.GetClients() {
		if !isExist {
			continue
		}
		GameEngine.AddPlayer(client.ID)
	}

	err := GameEngine.BuildPlayersCards()
	if err != nil {
		SendErrRes(c, Start, ErrForServerError, "伺服器 - 產生玩家牌組失敗")
		return
	}

	result := GameEngine.GetPlayersCards()
	if err != nil {
		BroadcastErrRes(c, Start, ErrForServerError, "伺服器 - 轉換玩家牌組資料失敗")
		return
	}

	BroadcastSuccessRes(c, UpdateAllDecks, result, "取得所有玩家當前牌組")
}

// msg_code 2
func OnHit(c *Client) {
	if !GameEngine.IsGameStart() {
		SendErrRes(c, Hit, ErrForNotCorrectState, "尚未開始遊戲 不可發牌  請某一方先開始遊戲")
		return
	}

	if GameEngine.IsPlayerEndTurn(c.ID) {
		SendErrRes(c, Hit, ErrForWrongFlow, "玩家不可再拿牌")
		return
	}

	err := GameEngine.DealCardToPlayer(c.ID)
	if err != nil {
		SendErrRes(c, Hit, ErrForDealCard, "發牌給玩家失敗")
		return
	}

	playerCardsPoint := GameEngine.CalculatePlayerCardsPoint(c.ID)
	if playerCardsPoint > MaxPoint {

	}

	result := GameEngine.GetPlayersCards()
	BroadcastSuccessRes(c, Hit, result, "取得所有玩家當前牌組")
}

// msg_code 3
func OnStand(c *Client) {
	if !GameEngine.IsGameStart() {
		SendErrRes(c, Stand, ErrForNotCorrectState, "尚未開始遊戲 不可發牌  請某一方先開始遊戲")
		return
	}

	GameEngine.StopPlayerAction(c.ID)
	BroadcastSuccessRes(c, Stand, c.ID, fmt.Sprintf("玩家ID-%s已經停牌", c.ID))
}
