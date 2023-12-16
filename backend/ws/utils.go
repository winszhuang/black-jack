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

func SendErrRes(c IClient, msgCode OperationCode, errorCode WSError, message string) {
	c.WsSend(WSResponse{
		Route:     PlayBlackJack,
		Data:      nil,
		Success:   false,
		ErrorCode: errorCode,
		Message:   message,
	}.Byte())
}

func SendSuccessRes(c IClient, msgCode OperationCode, data interface{}, message string) {
	c.WsSend(WSResponse{
		Route:     PlayBlackJack,
		Data:      data,
		Success:   true,
		ErrorCode: 0,
		Message:   message,
	}.Byte())
}

func BroadcastErrRes(game *Room, msgCode OperationCode, errorCode WSError, message string) {
	game.Broadcast(WSResponse{
		Route:     PlayBlackJack,
		Data:      nil,
		Success:   false,
		ErrorCode: errorCode,
		Message:   message,
	}.Byte())
}

func BroadcastSuccessRes(game *Room, msgCode OperationCode, data interface{}, message string) {
	game.Broadcast(WSResponse{
		Route:     PlayBlackJack,
		Data:      data,
		Success:   true,
		ErrorCode: 0,
		Message:   message,
	}.Byte())
}
