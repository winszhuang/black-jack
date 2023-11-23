package ws

func SendErrRes(c *Client, msgCode MsgCode, errorCode WSError, message string) {
	c.Write(WSResponse{
		MsgCode:   msgCode,
		Data:      nil,
		Success:   false,
		ErrorCode: errorCode,
		Message:   message,
	}.Byte())
}

func SendSuccessRes(c *Client, msgCode MsgCode, data interface{}, message string) {
	c.Write(WSResponse{
		MsgCode:   msgCode,
		Data:      data,
		Success:   true,
		ErrorCode: 0,
		Message:   message,
	}.Byte())
}

func BroadcastErrRes(c *Client, msgCode MsgCode, errorCode WSError, message string) {
	c.Game.Broadcast(WSResponse{
		MsgCode:   msgCode,
		Data:      nil,
		Success:   false,
		ErrorCode: errorCode,
		Message:   message,
	}.Byte())
}

func BroadcastSuccessRes(c *Client, msgCode MsgCode, data interface{}, message string) {
	c.Game.Broadcast(WSResponse{
		MsgCode:   msgCode,
		Data:      data,
		Success:   true,
		ErrorCode: 0,
		Message:   message,
	}.Byte())
}
