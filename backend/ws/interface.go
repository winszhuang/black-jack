package ws

import "black-jack/card"

type ClientDetail struct {
	ID    string    `json:"id"`
	Deck  card.Deck `json:"deck"`
	State UserState `json:"state"`
}

type IClient interface {
	InitPlayerInfo()
	GetID() string
	GetCurrentState() UserState
	UpdateCurrentState(state UserState)
	CalculateTotalPoints() int
	GetGameDetail() ClientDetail
	AddCard(card card.Card)
	WsSend([]byte)
	CloseWsSend()
}
