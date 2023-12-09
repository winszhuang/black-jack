package ws

import "black-jack/game"

type ClientDetail struct {
	ID    string    `json:"id"`
	Deck  game.Deck `json:"deck"`
	State UserState `json:"state"`
}

type IClient interface {
	InitPlayerInfo()
	GetID() string
	GetCurrentState() UserState
	UpdateCurrentState(state UserState)
	CalculateTotalPoints() int
	GetGameDetail() ClientDetail
	AddCard(card game.Card)
	WsSend([]byte)
	CloseWsSend()
}
