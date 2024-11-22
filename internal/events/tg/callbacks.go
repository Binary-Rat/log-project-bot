package tg

import (
	"context"
	"log"
	"strings"
)

func (p *Processor) doCallBack(text string, chatID int, username string) error {
	data := strings.TrimSpace(text)
	log.Printf("got new callback: %s from user: %s, chaID: %d\n", data, username, chatID)

	switch data {
	case atisuCall:
		return p.atisuCallback(username)
	case "a":
		return nil
	}
	return nil
}

func (p *Processor) atisuCallback(userID string) error {
	p.fsm.GetFilter(context.TODO(), userID)
	p.source.GetCarsWithFilter(nil)
	return nil
}
