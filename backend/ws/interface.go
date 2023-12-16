package ws

import (
	"black-jack/card"
	"github.com/google/uuid"
)

type ClientDetail struct {
	ID    uuid.UUID `json:"id"`
	Deck  card.Deck `json:"deck"`
	State UserState `json:"state"`
}

type IClient interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	InitPlayerInfo()
	GetID() uuid.UUID
	GetCurrentState() UserState
	UpdateCurrentState(state UserState)
	CalculateTotalPoints() int
	GetGameDetail() ClientDetail
	AddCard(card card.Card)
	WsSend([]byte)
	CloseWsSend()
	IsLogin() bool
	SetCurrRoom(*Room)
	GetCurrRoom() *Room
}
