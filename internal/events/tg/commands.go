package tg

import (
	"context"
	"fmt"
	"log"
	"log-proj/internal/text"
	"strconv"
	"strings"
)

const (
	//states
	stateCalcV = "calcV"
	stateCalcW = "calcW"
	//commands
	HelpCmd  = "/help"
	StartCmd = "/start"
	CalcCmd  = "/calc"
)

func (p *Processor) doCmd(msg string, chatID int, username string) error {
	msg = strings.TrimSpace(msg)
	log.Printf("get new command: %s from user: %s\n", msg, username)

	state := p.fsm.GetState(context.Background(), username)

	if msg == "/exit" {
		p.fsm.SetState(context.Background(), username, "")
		return p.tg.SendMessage(chatID, text.HelloMsg)
	}

	if state == stateCalcV {
		return p.calcVEvent(msg, chatID, username)
	}

	if state == stateCalcW {
		return p.calcWEvent(msg, chatID, username)
	}

	switch msg {
	case StartCmd:
		return p.sendHello(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case CalcCmd:
		return p.startCalc(chatID, username)
	default:
		return p.tg.SendMessage(chatID, text.MsgUnknownCommand)
	}
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, text.MsgUnknownCommand)
}

func (p *Processor) sendHello(chatID int, u string) error {
	return p.tg.SendMessage(chatID, fmt.Sprintf(text.HelloMsg, u))
}

func (p *Processor) startCalc(chatID int, u string) error {
	p.fsm.SetState(context.Background(), u, stateCalcV)
	return p.tg.SendMessage(chatID, text.CalcVMsg)
}

func (p *Processor) calcWEvent(msg string, chatID int, username string) error {
	w, err := strconv.ParseFloat(msg, 32)
	if err != nil {
		log.Printf("can`t convert string to int: %v", err)
		return p.tg.SendMessage(chatID, text.WrongValue)
	}
	p.fsm.SetLoadW(context.TODO(), username, w)
	v, m := p.fsm.GetLoad(context.TODO(), username)
	p.fsm.SetState(context.TODO(), username, "")
	log.Println(v, m)
	return p.tg.SendMessage(chatID, fmt.Sprintf("%f, %f", v, m))
}

func (p *Processor) calcVEvent(msg string, chatID int, username string) error {
	v, err := strconv.ParseFloat(msg, 64)
	if err != nil {
		log.Printf("can`t convert string to int: %v", err)
		return p.tg.SendMessage(chatID, text.WrongValue)
	}
	p.fsm.SetLoadV(context.TODO(), username, v)
	p.fsm.SetState(context.TODO(), username, stateCalcW)
	return p.tg.SendMessage(chatID, text.CalcWMsg)
}
