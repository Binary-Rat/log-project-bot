package tg

import (
	"errors"
	"log-proj/pkg/clients/tg"
	"log-proj/pkg/events"
	"log-proj/pkg/lib/e"
)

type Processor struct {
	tg     *tg.Client
	offset int
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnkownType = errors.New("UnkownType")
)

func New(client *tg.Client) *Processor {
	return &Processor{
		tg: client,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Warp("can`t get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Warp("can`t process event", ErrUnkownType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Warp("can`t get meta", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Warp("can`t do cmd", err)
	}
	return nil

}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Warp("can`t get meta", nil)
	}
	return res, nil
}

func event(upd tg.Updates) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd tg.Updates) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd tg.Updates) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
