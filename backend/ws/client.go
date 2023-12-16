package ws

import (
	"black-jack/card"
	"black-jack/utils"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Write pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the Room.
type Client struct {
	Center *GameCenter

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// user ID
	ID uuid.UUID `json:"id"`

	// user property
	property map[string]interface{}

	// user property lock
	propertyLock *sync.RWMutex

	// user card data
	playInfo *PlayInfo `json:"playInfo"`

	currRoom *Room
}

func (c *Client) IsLogin() bool {
	isLogin, err := c.GetProperty("isLogin")
	if err != nil {
		return false
	}
	return isLogin.(bool)
}

func NewClient(center *GameCenter, conn *websocket.Conn, ID uuid.UUID) *Client {
	return &Client{
		Center:       center,
		conn:         conn,
		send:         make(chan []byte, 256),
		ID:           ID,
		playInfo:     NewPlayInfo(),
		property:     map[string]interface{}{},
		propertyLock: &sync.RWMutex{},
	}
}

func NewPlayInfo() *PlayInfo {
	u := &PlayInfo{}
	u.Init()
	return u
}

func (u *PlayInfo) Init() {
	u.deck = make(card.Deck, 0)
	u.currentState = Wait
}

func (c *Client) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Client) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Client) InitPlayerInfo() {
	c.playInfo.Init()
}

func (c *Client) GetID() uuid.UUID {
	return c.ID
}

func (c *Client) GetCurrentState() UserState {
	return c.playInfo.currentState
}

func (c *Client) UpdateCurrentState(state UserState) {
	c.playInfo.currentState = state
}

func (c *Client) CalculateTotalPoints() int {
	return c.playInfo.deck.CalculateTotalPoints()
}

func (c *Client) GetGameDetail() ClientDetail {
	return ClientDetail{
		ID:    c.ID,
		Deck:  c.playInfo.deck,
		State: c.playInfo.currentState,
	}
}

func (c *Client) AddCard(card card.Card) {
	c.playInfo.deck = c.playInfo.deck.AddCard(card)
}

func (c *Client) WsSend(data []byte) {
	select {
	case c.send <- data:
	default:
		c.CloseWsSend()
		c.Center.RemoveClient(c)
	}
}

func (c *Client) CloseWsSend() {
	close(c.send)
}

func (c *Client) SetCurrRoom(room *Room) {
	c.currRoom = room
}

func (c *Client) GetCurrRoom() *Room {
	return c.currRoom
}

type PlayInfo struct {
	deck         card.Deck
	currentState UserState
}

type UserState int

const (
	_ UserState = iota
	Wait
	Ready
	Play
	Stop
	End
)

// readPump pumps messages from the websocket connection to the Room.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		if c.IsLogin() {
			c.Center.RemoveClient(c)
			c.currRoom.OnLeave(c)
		} else {
			c.Center.RemoveGuest(c)
		}
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var req WSRequest
		unmarshalErr := json.Unmarshal(message, &req)
		// #NOTICE 這居然檢測不出來!!
		if unmarshalErr != nil {
			c.Write([]byte("錯誤的回傳格式"))
			continue
		}

		switch req.Route {
		case Login:
			c.Center.HandleLogin(c, req.Data)
		case Register:
			c.Center.HandleRegister(c, req.Data)
		case JoinRoom:
			c.Center.HandleJoinRoom(c, req.Data)
		case LeaveRoom:
			c.Center.HandleLeaveRoom(c, req.Data)
		case PlayBlackJack:
			c.Center.HandlePlayBlackJack(c, req.Data)
		}
	}
}

func (c *Client) GetConn() *websocket.Conn {
	return c.conn
}

// writePump pumps messages from the Room to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The Room closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) Write(data []byte) {
	c.send <- data
}

// ServeWs handles websocket requests from the peer.
func ServeWs(center *GameCenter, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 玩家有遊戲中心 但遊戲中心暫不保存訪客資訊
	client := NewClient(center, conn, uuid.New())

	isLogin := checkUserLogin(r, client)
	if isLogin {
		center.AddClient(client)
		client.SetProperty("isLogin", true)
		// 發送所有房間資訊給玩家
		client.WsSend(GenSuccessRes(GetRoomsInfo, center.getRoomsInfo(), "所有房間資訊"))
	} else {
		center.AddGuest(client)
	}
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func checkUserLogin(r *http.Request, client *Client) bool {
	key, value, isExist := parseProtocol(r)
	if isExist && key == "access_token" && value != "" {
		_, err := utils.VerifyJWT(value)
		if err == nil {
			return true
			//userIdStr := claims["userId"].(string)
			//userId, _ := strconv.Atoi(userIdStr)
			//client.SetProperty(key, value)
			//client.SetProperty("userId", userId)
		}
	}
	return false
}

func parseProtocol(r *http.Request) (string, string, bool) {
	protocol := r.Header.Get("Sec-WebSocket-Protocol")
	if protocol == "" {
		return "", "", false
	}

	// protocol格式為: access_token, 789456123
	split := strings.Split(protocol, ",")
	if len(split) != 2 {
		return "", "", false
	}

	return strings.TrimSpace(split[0]), strings.TrimSpace(split[1]), true
}
