package tg

import (
	"errors"
	"log-proj/internal/events"
	"log-proj/internal/source"
	"log-proj/pkg/clients/tg"
	"log-proj/pkg/db"
	"log-proj/pkg/fsm"
	"log-proj/pkg/lib/e"
)

type Processor struct {
	tg               *tg.Client
	offset           int
	storage          db.Interface
	fsm              fsm.Interface
	processUnhandled bool
	Source           source.Interface
}

type MetaMessage struct {
	ChatID   int
	Username string
}

type MetaCallBack struct {
	QueryID  int
	ChatID   int
	Username string
}

var (
	ErrUnkownType = errors.New("UnkownType")
)

func New(client *tg.Client, storage db.Interface, fsm fsm.Interface, processUnhandled bool, source source.Interface) *Processor {
	return &Processor{
		tg:               client,
		storage:          storage,
		fsm:              fsm,
		processUnhandled: processUnhandled,
		Source:           source,
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

	if !p.processUnhandled {
		p.offset = updates[len(updates)-1].ID + 1
		p.processUnhandled = true
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
	case events.CallBack:
		return p.processCallBack(event)
	default:
		return e.Warp("can`t process event", ErrUnkownType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := metaMessage(event)
	if err != nil {
		return e.Warp("can`t get meta", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Warp("can`t do cmd", err)
	}
	return nil

}

func (p *Processor) processCallBack(event events.Event) error {
	meta, err := metaCallBack(event)
	if err != nil {
		return e.Warp("can`t get meta", err)
	}
	if err := p.doCallBack(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Warp("can`t do cmd", err)
	}
	return nil
}

func metaMessage(event events.Event) (MetaMessage, error) {
	res, ok := event.Meta.(MetaMessage)
	if !ok {
		return MetaMessage{}, e.Warp("can`t get meta", nil)
	}
	return res, nil
}

func metaCallBack(event events.Event) (MetaCallBack, error) {
	res, ok := event.Meta.(MetaCallBack)
	if !ok {
		return MetaCallBack{}, e.Warp("can`t get meta", nil)
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
		res.Meta = MetaMessage{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	} else if updType == events.CallBack {
		res.Meta = MetaCallBack{
			Username: upd.CallbackQuery.From.Username,
			ChatID:   upd.CallbackQuery.Message.Chat.ID,
		}
	}

	return res
}

func fetchText(upd tg.Updates) string {
	if upd.Message != nil {
		return upd.Message.Text
	} else if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}
	return ""
}

func fetchType(upd tg.Updates) events.Type {
	if upd.Message != nil {
		return events.Message
	} else if upd.CallbackQuery != nil {
		return events.CallBack
	}
	return events.Unknown
}
