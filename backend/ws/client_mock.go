package ws

import (
	"black-jack/game"
)

type ClientMock struct {
	ID       string
	playInfo *PlayInfo `json:"playInfo"`
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

func (c *ClientMock) AddCard(card game.Card) {
	c.playInfo.deck = c.playInfo.deck.AddCard(card)
}

func (c *ClientMock) WsSend(bytes []byte) {
	// nothing
}

func (c *ClientMock) CloseWsSend() {
	// nothing
}
