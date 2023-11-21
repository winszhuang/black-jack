package ws

import "encoding/json"

type MsgCode int

const (
	Null MsgCode = iota
	// 遊戲開始
	Start
	// 玩家要牌
	Hit
	// 玩家停牌
	Stand

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
