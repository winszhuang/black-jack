package ws

func GenErrRes(route Route, errorCode WSError, message string) []byte {
	return WSResponse{
		Route:     route,
		Data:      nil,
		Success:   false,
		ErrorCode: errorCode,
		Message:   message,
	}.Byte()
}

func GenErrResForInvalidRequest(route Route) []byte {
	return GenErrRes(route, ErrForInvalidRequest, "錯誤的參數")
}

func GenErrResForUnauthorized(route Route) []byte {
	return GenErrRes(route, ErrForUnauthorized, "未授權")
}

func GenSuccessRes(route Route, data interface{}, message string) []byte {
	return WSResponse{
		Route:     route,
		Data:      data,
		Success:   true,
		ErrorCode: 0,
		Message:   message,
	}.Byte()
}

func SendGameErrRes(c IClient, operationCode OperationCode, errorCode WSError, message string) {
	c.WsSend(WSResponse{
		Route:     PlayBlackJack,
		Success:   false,
		ErrorCode: errorCode,
		Message:   message,
		Data: WSPlayGameReqResData{
			OpCode:   operationCode,
			GameType: 1,
			GameData: nil,
		},
	}.Byte())
}

func SendGameSuccessRes(c IClient, operationCode OperationCode, data interface{}, message string) {
	c.WsSend(WSResponse{
		Route:     PlayBlackJack,
		Success:   true,
		ErrorCode: 0,
		Message:   message,
		Data: WSPlayGameReqResData{
			OpCode:   operationCode,
			GameType: 1,
			GameData: data,
		},
	}.Byte())
}

func BroadcastGameErrRes(game *Room, operationCode OperationCode, errorCode WSError, message string) {
	game.Broadcast(WSResponse{
		Route:     PlayBlackJack,
		Success:   false,
		ErrorCode: errorCode,
		Message:   message,
		Data: WSPlayGameReqResData{
			GameData: nil,
			OpCode:   operationCode,
			GameType: 1,
		},
	}.Byte())
}

func BroadcastGameSuccessRes(game *Room, operationCode OperationCode, data interface{}, message string) {
	game.Broadcast(WSResponse{
		Route:     PlayBlackJack,
		Success:   true,
		ErrorCode: 0,
		Message:   message,
		Data: WSPlayGameReqResData{
			GameData: data,
			OpCode:   operationCode,
			GameType: 1,
		},
	}.Byte())
}
