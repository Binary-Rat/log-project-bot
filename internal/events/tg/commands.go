package tg

import (
	"context"
	"fmt"
	"log"
	"log-proj/internal/text"
	"log-proj/pkg/clients/tg"
	"log-proj/pkg/lib/e"
	"log-proj/pkg/models"
	"strconv"
	"strings"

	"github.com/Binary-Rat/atisu"
)

const (
	//states
	stateCalcV = "calcV"
	stateCalcW = "calcW"
	//commands
	helpCmd  = "/help"
	startCmd = "/start"
	calcCmd  = "/calc"
	exitCmd  = "/exit"
)

//Unfortunately, i really dont know how to divede this logi(((
//The commands.go need to be in internal
//The tg.go need to be in pkg

func (p *Processor) doCmd(msg string, chatID int, username string) error {
	msg = strings.TrimSpace(msg)
	log.Printf("get new command: %s from user: %s\n", msg, username)

	if msg == exitCmd {
		return p.exit(username, chatID)
	}

	state := p.fsm.GetState(context.Background(), username)
	switch state {
	case stateCalcV:
		return p.calcVEvent(msg, chatID, username)
	case stateCalcW:
		return p.calcWEvent(msg, chatID, username)
	case stateCityFrom:
		return p.cityFromEvent(msg, chatID, username)
	case stateCityTo:
		return p.cityToEvent(msg, chatID, username)
	}

	//check command
	switch msg {
	case startCmd:
		return p.sendHello(chatID, username)
	case helpCmd:
		return p.sendHelp(chatID)
	case calcCmd:
		return p.startCalc(chatID, username)
	default:
		return p.tg.SendMessage(chatID, text.MsgUnknownCommand, nil)
	}
}

func (p *Processor) exit(username string, chatID int) error {
	p.fsm.SetState(context.TODO(), username, "")
	return p.tg.SendMessage(chatID, text.HelloMsg, nil)

}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, text.MsgUnknownCommand, nil)
}

func (p *Processor) sendHello(chatID int, u string) error {
	return p.tg.SendMessage(chatID, fmt.Sprintf(text.HelloMsg, u), nil)
}

func (p *Processor) startCalc(chatID int, u string) error {
	p.fsm.SetState(context.Background(), u, stateCalcV)
	return p.tg.SendMessage(chatID, text.CalcVMsg, nil)
}

func (p *Processor) calcWEvent(msg string, chatID int, username string) error {
	w, err := strconv.ParseFloat(msg, 32)
	if err != nil {
		log.Printf("can`t convert string to int: %v", err)
		return p.tg.SendMessage(chatID, text.WrongValue, nil)
	}
	p.fsm.SetLoadW(context.TODO(), username, w)
	p.fsm.SetState(context.TODO(), username, "")

	cars, err := p.cars(username)
	if err != nil {
		e.Warp("can`t get cars", err)
	}
	keybord := &tg.ReplyMarkup{InlineKeyboard: keybord}
	return p.tg.SendMessage(chatID, strings.Join(cars.Names(), " "), keybord)
}

func (p *Processor) calcVEvent(msg string, chatID int, username string) error {
	v, err := strconv.ParseFloat(msg, 64)
	if err != nil {
		log.Printf("can`t convert string to int: %v", err)
		return p.tg.SendMessage(chatID, text.WrongValue, nil)
	}
	p.fsm.SetLoadV(context.TODO(), username, v)
	p.fsm.SetState(context.TODO(), username, stateCalcW)
	return p.tg.SendMessage(chatID, text.CalcWMsg, nil)
}

func (p *Processor) cityFromEvent(msg string, chatID int, username string) error {
	p.fsm.SetCityFrom(context.TODO(), username, msg)
	p.fsm.SetState(context.TODO(), username, stateCityTo)
	return p.tg.SendMessage(chatID, "Введите город назначения", nil)
}

func (p *Processor) cityToEvent(msg string, chatID int, username string) error {
	p.fsm.SetCityTO(context.TODO(), username, msg)
	log.Println("searching for cars")
	data := p.fsm.GetRoadCities(context.TODO(), username)
	cities, err := p.source.GetCityID(data)
	if err != nil {
		return fmt.Errorf("can`t get city id: %w", err)
	}
	filter := atisu.Filter{}
	filter.From.ID = (*cities)[data[0]].CityID
	filter.To.ID = (*cities)[data[1]].CityID
	log.Println(filter)
	cars, err := p.source.GetCarsWithFilter(filter)
	if err != nil {
		return fmt.Errorf("can`t get cars: %w", err)
	}
	return p.tg.SendMessage(chatID, strings.Join(cars.Names(), " "), nil)
}

func (p *Processor) cars(userID string) (cars models.Cars, err error) {
	v, w := p.fsm.GetLoad(context.TODO(), userID)
	return p.storage.GetCars(v, w)
}
