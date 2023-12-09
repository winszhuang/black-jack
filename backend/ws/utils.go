package ws

func SendErrRes(c IClient, msgCode MsgCode, errorCode WSError, message string) {
	c.WsSend(WSResponse{
		MsgCode:   msgCode,
		Data:      nil,
		Success:   false,
		ErrorCode: errorCode,
		Message:   message,
	}.Byte())
}

func SendSuccessRes(c IClient, msgCode MsgCode, data interface{}, message string) {
	c.WsSend(WSResponse{
		MsgCode:   msgCode,
		Data:      data,
		Success:   true,
		ErrorCode: 0,
		Message:   message,
	}.Byte())
}

func BroadcastErrRes(game *Game, msgCode MsgCode, errorCode WSError, message string) {
	game.Broadcast(WSResponse{
		MsgCode:   msgCode,
		Data:      nil,
		Success:   false,
		ErrorCode: errorCode,
		Message:   message,
	}.Byte())
}

func BroadcastSuccessRes(game *Game, msgCode MsgCode, data interface{}, message string) {
	game.Broadcast(WSResponse{
		MsgCode:   msgCode,
		Data:      data,
		Success:   true,
		ErrorCode: 0,
		Message:   message,
	}.Byte())
}
