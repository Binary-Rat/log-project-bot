package tg

import (
	"context"
	"log"
	"strings"
)

const (
	stateCityFrom = "cityFrom"
	stateCityTo   = "cityTo"
)

func (p *Processor) doCallBack(text string, chatID int, username string) error {
	data := strings.TrimSpace(text)
	log.Printf("got new callback: %s from user: %s, chaID: %d\n", data, username, chatID)

	switch data {
	case atisuCall:
		return p.atisuCallback(username, chatID)
	case "a":
		return nil
	}
	return nil
}

func (p *Processor) atisuCallback(userID string, chatID int) error {
	p.fsm.SetState(context.TODO(), userID, stateCityFrom)
	return p.tg.SendMessage(chatID, "Укажите город отправления", nil)
}
