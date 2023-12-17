package ws

import (
	"encoding/json"
	"github.com/google/uuid"
)

type Route string

const (
	Register      Route = "register"
	Login         Route = "login"
	JoinRoom      Route = "join_room"
	LeaveRoom     Route = "leave_room"
	PlayBlackJack Route = "play_black_jack"
	GetRoomsInfo  Route = "get_rooms_info"
	WsConnected   Route = "ws_connected"
)

type OperationCode int

const (
	_ OperationCode = iota

	// 玩家加入
	ClientJoin

	// 廣播某玩家進入遊戲
	BroadcastJoin

	// 廣播某玩家離開遊戲
	BroadcastLeave

	// 玩家準備
	ClientReady

	// 廣播某玩家按下準備
	BroadcastReady

	// 廣播遊戲開始
	BroadcastGameStart

	// 玩家要牌
	ClientHit

	// 廣播某玩家要牌
	BroadcastHit

	// 廣播有人爆牌
	BroadcastBust

	// 玩家停止要牌
	ClientStand

	// 廣播某玩家停止要牌
	BroadcastStand

	// 廣播遊戲結束
	BroadcastGameOver

	// 廣播遊戲重新開始
	BroadcastReStart

	// 更新所有玩家資訊
	UpdatePlayersDetail
)

type WSRequest struct {
	Route Route       `json:"route"`
	Data  interface{} `json:"data"`
}

type WSLoginReqData struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type WSRegisterReqData struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type WSJoinRoomReqData struct {
	RoomID uuid.UUID `json:"room_id"`
}

type WSLeaveRoomReqData struct {
	RoomID uuid.UUID `json:"room_id"`
}

type WSPlayGameReqData struct {
	OpCode   OperationCode `json:"opcode"`
	GameType int           `json:"gametype"`
	GameData interface{}   `json:"gamedata"`
}

type WSResponse struct {
	Route     Route       `json:"route"`
	Data      interface{} `json:"data"`
	Success   bool        `json:"success" example:"true"`
	ErrorCode WSError     `json:"error_code" example:"0"`
	Message   string      `json:"message" example:"回傳成功"`
}

func (r WSResponse) Byte() []byte {
	marshal, _ := json.Marshal(r)
	return marshal
}
