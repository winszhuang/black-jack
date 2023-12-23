package ws

import "C"
import (
	"black-jack/models"
	"black-jack/repository"
	"black-jack/utils"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"sync"
)

type GameCenter struct {
	rooms    []*Room          // 所有房間
	clients  map[IClient]bool // 註冊的所有玩家
	guests   map[IClient]bool // 訪客
	mu       *sync.RWMutex    // 鎖
	userRepo repository.IUserRepository
}

func NewGameCenter(userRepo repository.IUserRepository, multiGame ...*Room) *GameCenter {
	rooms := make([]*Room, 0)
	if len(multiGame) > 0 {
		for _, game := range multiGame {
			rooms = append(rooms, game)
		}
	}
	return &GameCenter{
		rooms:    rooms,
		clients:  make(map[IClient]bool),
		guests:   make(map[IClient]bool),
		userRepo: userRepo,
		mu:       &sync.RWMutex{},
	}
}

func (gc *GameCenter) AddClient(c IClient) {
	gc.mu.Lock()
	if _, isExist := gc.clients[c]; !isExist {
		gc.clients[c] = true
	}
	gc.mu.Unlock()
}

func (gc *GameCenter) RemoveClient(c IClient) {
	gc.mu.Lock()
	if _, isExist := gc.clients[c]; isExist {
		delete(gc.clients, c)
	}
	gc.mu.Unlock()
}

func (gc *GameCenter) AddGuest(c *Client) {
	gc.mu.Lock()
	if _, isExist := gc.guests[c]; !isExist {
		gc.guests[c] = true
	}
	gc.mu.Unlock()
}

func (gc *GameCenter) RemoveGuest(c IClient) {
	gc.mu.Lock()
	if _, isExist := gc.guests[c]; isExist {
		delete(gc.guests, c)
	}
	gc.mu.Unlock()
}

type LoginResponse struct {
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"user_id"`
}

func (gc *GameCenter) HandleLogin(c IClient, data interface{}) {
	if c.IsLogin() {
		c.WsSend(GenErrRes(Login, ErrForDuplicateLogin, "已經登入過了!!"))
		return
	}

	mapData, ok := data.(map[string]interface{})
	if !ok {
		c.WsSend(GenErrResForInvalidRequest(Login))
		return
	}

	var req WSLoginReqData
	if mapstructure.Decode(mapData, &req) != nil {
		c.WsSend(GenErrResForInvalidRequest(Login))
		return
	}

	user, err := gc.userRepo.GetByName(req.UserName)
	if err != nil {
		c.WsSend(GenErrRes(Login, 0, "無此使用者"))
		return
	}

	err = utils.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.WsSend(GenErrRes(Login, 0, "錯誤的密碼!!"))
		return
	}

	// generate jwt
	jwtStr, err := utils.CreateJWT(map[string]string{
		"userId":   user.ID.String(),
		"username": user.Name,
	})
	utils.AssertNoError(err)

	c.WsSend(GenSuccessRes(Login, LoginResponse{
		Token:  jwtStr,
		UserID: user.ID,
	}, "登入成功"))

	// 登入後加入
	gc.AddClient(c)
	c.SetProperty("isLogin", true)
	c.UpdateLoginInfo(&LoginInfo{
		UserName: user.Name,
		UserID:   user.ID.String(),
	})

	// 發送所有房間資訊給玩家
	c.WsSend(GenSuccessRes(GetRoomsInfo, gc.getRoomsInfo(), "所有房間資訊"))
}

type RoomInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (gc *GameCenter) getRoomsInfo() []RoomInfo {
	var roomsInfo []RoomInfo
	for _, room := range gc.rooms {
		roomsInfo = append(roomsInfo, RoomInfo{
			ID:   room.ID,
			Name: room.Name,
		})
	}
	return roomsInfo
}

func (gc *GameCenter) HandleRegister(c IClient, data interface{}) {
	mapData, ok := data.(map[string]interface{})
	if !ok {
		c.WsSend(GenErrResForInvalidRequest(Register))
		return
	}

	var req WSRegisterReqData
	if mapstructure.Decode(mapData, &req) != nil {
		c.WsSend(GenErrResForInvalidRequest(Register))
		return
	}

	_, err := gc.userRepo.GetByName(req.UserName)
	if err == nil {
		c.WsSend(GenErrRes(Register, 0, "此暱稱已被使用 請重新選擇"))
		return
	}

	hashPassword, err := utils.GenerateFromPassword([]byte(req.Password))
	utils.AssertNoError(err)

	var user models.User
	user.Name = req.UserName
	user.Password = string(hashPassword)
	if gc.userRepo.Create(&user) != nil {
		c.WsSend(GenErrRes(Register, 0, "創建使用者失敗"))
		return
	}

	c.WsSend(GenSuccessRes(Register, true, "註冊成功"))
}

func (gc *GameCenter) HandleJoinRoom(c IClient, data interface{}) {
	if !c.IsLogin() {
		c.WsSend(GenErrResForUnauthorized(JoinRoom))
		return
	}

	mapData, ok := data.(map[string]interface{})
	if !ok {
		c.WsSend(GenErrResForInvalidRequest(JoinRoom))
		return
	}

	roomUUID, isOk := getUUIDByMapData(mapData, "room_id")
	if !isOk {
		c.WsSend(GenErrResForInvalidRequest(JoinRoom))
		return
	}

	room := gc.findRoomByID(roomUUID)
	if room == nil {
		c.WsSend(GenErrRes(JoinRoom, ErrForCantFindThisRoom, "無此房間"))
		return
	}

	c.WsSend(GenSuccessRes(JoinRoom, room.ID, "進入房間"))

	room.OnJoin(c)
	c.SetCurrRoom(room)
}

func (gc *GameCenter) HandleLeaveRoom(c IClient, data interface{}) {
	if !c.IsLogin() {
		c.WsSend(GenErrResForUnauthorized(LeaveRoom))
		return
	}

	mapData, ok := data.(map[string]interface{})
	if !ok {
		c.WsSend(GenErrResForInvalidRequest(LeaveRoom))
		return
	}

	roomUUID, isOk := getUUIDByMapData(mapData, "room_id")
	if !isOk {
		c.WsSend(GenErrResForInvalidRequest(JoinRoom))
		return
	}

	room := gc.findRoomByID(roomUUID)
	if room == nil {
		c.WsSend(GenErrRes(LeaveRoom, ErrForCantFindThisRoom, "無此房間"))
		return
	}

	room.OnLeave(c)
	c.SetCurrRoom(nil)
}

func (gc *GameCenter) HandlePlayBlackJack(c IClient, data interface{}) {
	if !c.IsLogin() {
		c.WsSend(GenErrResForUnauthorized(PlayBlackJack))
		return
	}

	mapData, ok := data.(map[string]interface{})
	if !ok {
		c.WsSend(GenErrResForInvalidRequest(PlayBlackJack))
		return
	}

	var req WSPlayGameReqResData
	if mapstructure.Decode(mapData, &req) != nil {
		c.WsSend(GenErrResForInvalidRequest(PlayBlackJack))
		return
	}

	currRoom := c.GetCurrRoom()
	if currRoom == nil {
		c.WsSend(GenErrRes(PlayBlackJack, ErrForClientNotInRoom, "玩家不在任何房間"))
		return
	}

	switch req.OpCode {
	case ClientReady:
		currRoom.OnReady(c)
	case ClientHit:
		currRoom.OnHit(c)
	case ClientStand:
		currRoom.OnStand(c)
	}
}

func (gc *GameCenter) findRoomByID(id uuid.UUID) *Room {
	for _, room := range gc.rooms {
		if room.ID == id {
			return room
		}
	}
	return nil
}

func getUUIDByMapData(mapData map[string]interface{}, key string) (uuid.UUID, bool) {
	roomId, isExist := mapData[key]
	if !isExist {
		return uuid.UUID{}, false
	}
	roomIdStr, isConvertOk := roomId.(string)
	if !isConvertOk {
		return uuid.UUID{}, false
	}

	roomUUID, err := uuid.Parse(roomIdStr)
	if err != nil {
		return uuid.UUID{}, false
	}
	return roomUUID, true
}
