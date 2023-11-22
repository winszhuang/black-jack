package ws

import (
	"black-jack/game"
	"math/rand"
)

var GameEngine *game.Game

func init() {
	cardDealer := game.NewCardDealer(rand.New(rand.NewSource(848484814486486)))
	GameEngine = game.NewGame(cardDealer)
}

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

func OnHit(c *Client) {
	if !GameEngine.IsGameStart() {
		SendErrRes(c, Hit, ErrForNotCorrectState, "尚未開始遊戲 不可發牌  請某一方先開始遊戲")
		return
	}

	err := GameEngine.DealCardToPlayer(c.ID)
	if err != nil {
		SendErrRes(c, Hit, ErrForDealCard, "發牌給玩家失敗")
		return
	}

	result := GameEngine.GetPlayersCards()
	if err != nil {
		BroadcastErrRes(c, Hit, ErrForGetAllPlayerCards, "伺服器 - 轉換玩家牌組資料失敗")
		return
	}

	BroadcastSuccessRes(c, Hit, result, "取得所有玩家當前牌組")
}

func OnStand(c *Client) {
	if !GameEngine.IsGameStart() {
		SendErrRes(c, Stand, ErrForNotCorrectState, "尚未開始遊戲 不可發牌  請某一方先開始遊戲")
		return
	}
}
