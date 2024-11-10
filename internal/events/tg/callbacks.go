package tg

import (
	"log"
	"strings"
)

func (p *Processor) doCallBack(text string, chatID int, username string) error {
	data := strings.TrimSpace(text)
	log.Printf("got new callback: %s from user: %s, chaID: %d\n", data, username, chatID)
	return nil
}
