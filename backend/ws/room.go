package ws

import (
	"black-jack/card"
	"fmt"
	"github.com/google/uuid"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"log"
	"sync"
	"time"
)

const MaxPoint = 21
const DealerClientID = "dealer"

type Room struct {
	ID uuid.UUID

	Name string // 房間名稱

	clients *orderedmap.OrderedMap[IClient, bool] // 註冊的所有玩家

	cardDealer card.ICardDealer // 發牌員

	mu *sync.RWMutex // 鎖
}

func NewRoom(name string, cardDealer card.ICardDealer) *Room {
	g := &Room{
		clients:    orderedmap.New[IClient, bool](),
		cardDealer: cardDealer,
		mu:         &sync.RWMutex{},
		ID:         uuid.New(),
		Name:       name,
	}
	g.Restart()
	return g
}

// Restart 重新開始遊戲
func (r *Room) Restart() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cardDealer.InitializeDeck()
	r.cardDealer.ShuffleDeck()
	for pair := r.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		client.InitPlayerInfo()
	}
}

func (r *Room) checkAllPlayerReadyToStart() {
	if !r.isAllPlayerSameState(Ready) {
		return
	}

	go func() {
		// 等待1秒，開始遊戲
		time.Sleep(time.Second)
		// 開始遊戲
		r.onGameStart()
	}()
}

func (r *Room) checkPlayerBustThenStop(c IClient) {
	if r.isPlayerBust(c) {
		r.mu.Lock()
		c.UpdateCurrentState(Stop)
		r.mu.Unlock()

		BroadcastSuccessRes(r, BroadcastStand, c.GetID(), fmt.Sprintf("ClientID-%s玩家已經停止動作", c.GetID()))
		BroadcastSuccessRes(r, UpdatePlayersDetail, r.getAllClientDetail(), "更新所有玩家資料")
	}
}

func (r *Room) checkAllPlayerStopThenEnd() {
	if !r.isAllPlayerSameState(Stop) {
		return
	}

	go func() {
		// 等待1秒，準備結束遊戲
		time.Sleep(time.Second * 1)
		// 結束遊戲
		r.onGameEnd()
	}()
}

// 玩家是否爆牌
func (r *Room) isPlayerBust(c IClient) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	totalPoints := c.CalculateTotalPoints()
	return totalPoints > MaxPoint
}

func (r *Room) updateAllPlayerState(state UserState) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for pair := r.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		client.UpdateCurrentState(state)
	}
}

func (r *Room) calculateFinalWinners() ([]IClient, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 玩家沒有1人以上就不用算了
	if r.clients.Len() < 2 {
		return []IClient{}, false
	}

	// 紀錄玩家得分-玩家資訊
	pointClientMap := map[int][]IClient{}
	for pair := r.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		clientPoint := client.CalculateTotalPoints()
		// 爆點不理他
		if clientPoint > MaxPoint {
			continue
		}

		if _, isExist := pointClientMap[clientPoint]; !isExist {
			pointClientMap[clientPoint] = []IClient{client}
		} else {
			pointClientMap[clientPoint] = append(pointClientMap[clientPoint], client)
		}
	}

	// 算出最大點數
	maxPoint := 0
	for point, _ := range pointClientMap {
		if point > maxPoint {
			maxPoint = point
		}
	}
	if _, isExist := pointClientMap[maxPoint]; isExist {
		return pointClientMap[maxPoint], true
	}
	log.Println("靠北 這邊也會錯哦 算分", maxPoint)
	return []IClient{}, false
}

func (r *Room) Broadcast(data []byte) {
	r.mu.Lock()
	for pair := r.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		client.WsSend(data)
	}
	r.mu.Unlock()
}

func (r *Room) getAllClientDetail() []ClientDetail {
	r.mu.Lock()
	defer r.mu.Unlock()

	var allClientsDetail []ClientDetail
	for pair := r.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key

		allClientsDetail = append(allClientsDetail, client.GetGameDetail())
	}
	return allClientsDetail
}

type ReadyNotification struct {
	ID string `json:"id"`
}

func (r *Room) isAllPlayerSameState(state UserState) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	count := 0
	for pair := r.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		if client.GetCurrentState() == state {
			count++
		}
	}
	return count == r.clients.Len()
}

func (r *Room) buildAllPlayerCards() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for pair := r.clients.Oldest(); pair != nil; pair = pair.Next() {
		client := pair.Key
		card1, err := r.cardDealer.DealCard()
		card2, err2 := r.cardDealer.DealCard()
		if err != nil {
			return err
		}
		if err2 != nil {
			return err2
		}
		client.AddCard(card1)
		client.AddCard(card2)
	}
	return nil
}

func (r *Room) isMoreThanOnePlayer() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.clients.Len() > 1
}
