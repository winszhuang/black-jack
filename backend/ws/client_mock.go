package ws

import (
	"black-jack/card"
)

type ClientMock struct {
	ID       string
	playInfo *PlayInfo `json:"playInfo"`
}

func (c *ClientMock) IsLogin() bool {
	return true
}

func (c *ClientMock) SetCurrRoom(room *Room) {
	// nothing
}

func (c *ClientMock) GetCurrRoom() *Room {
	return nil
}

func (c *ClientMock) SetProperty(key string, value interface{}) {
	// nothing
}

func (c *ClientMock) GetProperty(key string) (interface{}, error) {
	// nothing
	return nil, nil
}

func NewClientMock(ID string) *ClientMock {
	return &ClientMock{ID: ID, playInfo: NewPlayInfo()}
}

func (c *ClientMock) InitPlayerInfo() {
	c.playInfo.Init()
}

func (c *ClientMock) GetID() string {
	return c.ID
}

func (c *ClientMock) GetCurrentState() UserState {
	return Play
}

func (c *ClientMock) UpdateCurrentState(state UserState) {
	// nothing
}

func (c *ClientMock) CalculateTotalPoints() int {
	return c.playInfo.deck.CalculateTotalPoints()
}

func (c *ClientMock) GetGameDetail() ClientDetail {
	return ClientDetail{}
}

func (c *ClientMock) AddCard(card card.Card) {
	c.playInfo.deck = c.playInfo.deck.AddCard(card)
}

func (c *ClientMock) WsSend(bytes []byte) {
	// nothing
}

func (c *ClientMock) CloseWsSend() {
	// nothing
}
