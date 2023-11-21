package handler

import (
	"black-jack/game"
	"black-jack/ws"
	"encoding/json"
	"math/rand"
)

var GameEngine *game.Game

func init() {
	cardDealer := game.NewCardDealer(rand.New(rand.NewSource(848484814486486)))
	GameEngine = game.NewGame(cardDealer)
}

func OnStart(c *ws.Client) {
	if GameEngine.IsGameStart() {
		c.Write(ws.WSResponse{
			MsgCode:   ws.Start,
			Data:      nil,
			Success:   false,
			ErrorCode: 0,
			Message:   "遊戲已開始 不得重新開始遊戲",
		}.Byte())
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
		c.Hub.Broadcast(ws.WSResponse{
			MsgCode:   ws.Start,
			Data:      nil,
			Success:   false,
			ErrorCode: 0,
			Message:   "伺服器 - 產生玩家牌組失敗",
		}.Byte())
		return
	}

	result := GameEngine.GetPlayersCards()
	marshal, err := json.Marshal(result)
	if err != nil {
		c.Hub.Broadcast(ws.WSResponse{
			MsgCode:   ws.Start,
			Data:      nil,
			Success:   false,
			ErrorCode: 0,
			Message:   "伺服器 - 轉換玩家牌組資料失敗",
		}.Byte())
		return
	}

	c.Hub.Broadcast(ws.WSResponse{
		MsgCode:   ws.UpdateAllDecks,
		Data:      marshal,
		Success:   true,
		ErrorCode: 0,
		Message:   "取得所有玩家當前牌組",
	}.Byte())
}

func OnHit(c *ws.Client) {
	err := GameEngine.DealCardToPlayer(c.ID)
	if err != nil {
		c.Write(ws.WSResponse{
			MsgCode:   0,
			Data:      nil,
			Success:   false,
			ErrorCode: 0,
			Message:   "發牌給玩家失敗",
		}.Byte())
		return
	}

	result := GameEngine.GetPlayersCards()
	marshal, err := json.Marshal(result)
	if err != nil {
		c.Hub.Broadcast(ws.WSResponse{
			MsgCode:   0,
			Data:      nil,
			Success:   false,
			ErrorCode: 0,
			Message:   "伺服器 - 轉換玩家牌組資料失敗",
		}.Byte())
		return
	}

	c.Hub.Broadcast(ws.WSResponse{
		MsgCode:   ws.UpdateAllDecks,
		Data:      marshal,
		Success:   true,
		ErrorCode: 0,
		Message:   "取得所有玩家當前牌組",
	}.Byte())
}

func OnStand(c *ws.Client) {

}
