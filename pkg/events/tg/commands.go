package tg

import (
	"log"
	"strings"
)

const (
	HelpCmd           = "/help"
	StartCmd          = "/start"
	CalcCmd           = "/calc"
	msgUnknownCommand = "unknown command"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("get new command: %s from user: %s\n", text)
	switch text {
	case HelpCmd:
		return p.sendHelp(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgUnknownCommand)
}
