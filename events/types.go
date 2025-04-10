package events

import (
	"SPBHistoryBot/lib/e"
	"errors"
)

var (
	ErrUnknownMetaType = errors.New("unknown meta type")
)

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(event Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

type Meta struct {
	ChatId   int
	Username string
}

func (event Event) GetMeta() (Meta, error) {

	res, ok := event.Meta.(*Meta)

	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return *res, nil
}
