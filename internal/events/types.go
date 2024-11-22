package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(events Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
	CallBack
)

type Event struct {
	Type Type
	Text string
	Meta any
}
