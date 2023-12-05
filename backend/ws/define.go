package ws

import "encoding/json"

type MsgCode int

const (
	_ MsgCode = iota
	OneJoin
	// 某玩家進入遊戲
	BroadcastJoin
	// 某玩家離開遊戲
	BroadcastLeave

	OneReady
	// 某玩家按下準備
	BroadcastReady
	// 遊戲開始
	BroadcastGameStart

	OneHit
	// 玩家要牌
	BroadcastHit

	OneStand

	// 玩家停止要牌
	BroadcastStand
	// 遊戲結束
	BroadcastGameOver

	// 更新所有玩家資訊
	UpdatePlayersDetail
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
