package ws

import "encoding/json"

type MsgCode int

const (
	Null MsgCode = iota
	// 某玩家進入遊戲
	SomeOneJoin
	// 某玩家離開遊戲
	SomeOneLeave
	// 某玩家按下準備
	SomeOneReady
	// 遊戲開始
	GameStart
	// 玩家要牌
	SomeOneHit
	// 玩家停止要牌
	SomeOneStand

	GameOver

	UpdatePlayersDetail

	UpdateAllDecks
)

type WSRequest struct {
	Data    interface{} `json:"data"`
	MsgCode MsgCode     `json:"msg_code"`
}

type WSResponse struct {
	MsgCode   MsgCode     `json:"msg_code"`
	Data      interface{} `json:"data"`
	Success   bool        `json:"success" example:"true"`
	ErrorCode WSError     `json:"error_code" example:"0"`
	Message   string      `json:"message" example:"回傳成功"`
}

func (r WSResponse) Byte() []byte {
	marshal, _ := json.Marshal(r)
	return marshal
}
