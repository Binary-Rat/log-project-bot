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
	return p.tg.SendMessage(chatID, fmt.Sprintf(text.HelloMsg, username), nil)

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
	keybord := &tg.ReplyMarkup{InlineKeyboard: atiSearchKeyboard}
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
	// Используем контекст с тайм-аутом или передаем его как аргумент
	ctx := context.Background()
	p.fsm.SetState(context.TODO(), username, "")

	// Сохраняем город для обработки
	p.fsm.SetCityTO(ctx, username, msg)
	log.Println("searching for cars")

	// Получаем список городов
	data := p.fsm.GetRoadCities(ctx, username)
	cities, err := p.source.GetCityID(data)
	if err != nil {
		return fmt.Errorf("can't get city id for cities: %v, error: %w", data, err)
	}

	// Проверяем наличие городов и их идентификаторов
	if len(*cities) < 2 {
		return fmt.Errorf("insufficient city data: %v", *cities)
	}

	// Извлекаем ID городов
	fromCityID := (*cities)[data[0]].CityID
	toCityID := (*cities)[data[1]].CityID

	// Создание фильтра
	filter := atisu.Filter{
		Dates: atisu.DateOption{DateOption: "today"},
		From: atisu.CityFilter{
			ID:   fromCityID,
			Type: 2,
		},
		To: atisu.CityFilter{
			ID:   toCityID,
			Type: 2,
		},
	}
	// Получаем минимальные и максимальные значения для объема и веса
	filter.Volume.Min, filter.Weight.Min = p.fsm.GetLoad(ctx, username)
	filter.Volume.Max, filter.Weight.Max = filter.Volume.Min+100, filter.Weight.Min+100

	// Получаем машины с фильтром
	cars, err := p.source.GetCarsWithFilter(filter)
	if err != nil {
		return fmt.Errorf("can't get cars with filter %+v: %w", filter, err)
	}

	// Отправка сообщения
	if len(cars.Names()) == 0 {
		return p.tg.SendMessage(chatID, "No cars found", nil)
	}

	return p.tg.SendMessage(chatID, strings.Join(cars.Names(), " "), nil)
}

func (p *Processor) cars(userID string) (cars models.Cars, err error) {
	v, w := p.fsm.GetLoad(context.TODO(), userID)
	return p.storage.GetCars(v, w)
}
