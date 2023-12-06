package ws

import (
	"black-jack/game"
	"black-jack/utils"
	"encoding/json"
	"log"
	"net/http"
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

// Client is a middleman between the websocket connection and the Game.
type Client struct {
	Game *Game

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// user ID
	ID string

	// user game data
	playInfo *PlayInfo
}

func NewClient(game *Game, conn *websocket.Conn, ID string) *Client {
	return &Client{
		Game:     game,
		conn:     conn,
		send:     make(chan []byte, 256),
		ID:       ID,
		playInfo: NewPlayInfo(),
	}
}

type PlayInfo struct {
	deck         game.Deck
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

func NewPlayInfo() *PlayInfo {
	u := &PlayInfo{}
	u.Init()
	return u
}

func (u *PlayInfo) Init() {
	u.deck = make(game.Deck, 0)
	u.currentState = Wait
}

// readPump pumps messages from the websocket connection to the Game.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.Game.OnLeave(c)
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
		if req.MsgCode == 0 {
			c.Write(WSResponse{
				MsgCode:   0,
				Data:      nil,
				Success:   false,
				ErrorCode: ErrForWrongRequestFormat,
				Message:   "msgCode必須大於0",
			}.Byte())
			continue
		}

		switch req.MsgCode {
		case ClientReady:
			c.Game.OnReady(c)
		case ClientHit:
			c.Game.OnHit(c)
		case ClientStand:
			c.Game.OnStand(c)
		}
	}
}

func (c *Client) GetConn() *websocket.Conn {
	return c.conn
}

// writePump pumps messages from the Game to the websocket connection.
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
				// The Game closed the channel.
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
func ServeWs(game *Game, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(game, conn, utils.RandomPlayerName())
	client.Game.OnJoin(client)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
