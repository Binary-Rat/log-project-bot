package tg

import (
	"context"
	"log"
	"log-proj/internal/text"
	"strings"
)

const (
	HelpCmd           = "/help"
	StartCmd          = "/start"
	CalcCmd           = "/calc"
	msgUnknownCommand = "unknown command"
)

func (p *Processor) doCmd(msg string, chatID int, username string) error {
	msg = strings.TrimSpace(msg)
	log.Printf("get new command: %s from user: %s\n", msg, username)

	switch msg {
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case CalcCmd:
		p.fsm.SetState(context.Background(), username, "calcV")
		return p.tg.SendMessage(chatID, text.CalcVMsg)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgUnknownCommand)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, text.HelloMsg)
}
