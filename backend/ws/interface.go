package ws

import (
	"black-jack/card"
)

type ClientDetail struct {
	ID    string    `json:"id"`
	Deck  card.Deck `json:"deck"`
	State UserState `json:"state"`
	Name  string    `json:"name"`
}

type LoginInfo struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
}

type IClient interface {
	UpdateLoginInfo(*LoginInfo)
	GetLoginInfo() *LoginInfo
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	InitPlayerInfo()
	GetID() string
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
