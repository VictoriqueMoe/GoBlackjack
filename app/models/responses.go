package models

import (
	"github.com/gofiber/fiber/v2/log"
)

type ResponseMsg struct {
	Device      string     `json:"device"`
	Cards       []string   `json:"cards"`
	DealerCards []string   `json:"dealerCards"`
	HandValue   int        `json:"handValue"`
	DealerValue int        `json:"dealerValue"`
	Status      PlayStatus `json:"status"`
}

type ErrorMsg struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type StatusMsg struct {
	Wins  int `json:"wins"`
	Loses int `json:"loses"`
	Draws int `json:"draws"`
}

func NewStatusMsg(wins int, loses int, draws int) *StatusMsg {
	return &StatusMsg{Wins: wins, Loses: loses, Draws: draws}
}

func NewErrorMsg(msg string, err error, status int) *ErrorMsg {
	if err != nil {
		log.Error(err)
	}
	return &ErrorMsg{Status: status, Message: msg}
}

func NewResponseMsg(
	device string,
	cards []string,
	dealerCards []string,
	handValue int,
	dealerValue int,
	status PlayStatus,
) *ResponseMsg {
	return &ResponseMsg{
		Device:      device,
		Cards:       cards,
		DealerCards: dealerCards,
		HandValue:   handValue,
		DealerValue: dealerValue,
		Status:      status,
	}
}
