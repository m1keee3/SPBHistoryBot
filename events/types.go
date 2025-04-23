package events

import (
	"SPBHistoryBot/lib/e"
	"encoding/json"
	"errors"
	"log"
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
	Location
	Callback
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

type Meta struct {
	ChatID    int
	Username  string
	MessageID int
	Latitude  float64
	Longitude float64
}

func (event Event) GetMeta() (Meta, error) {

	res, ok := event.Meta.(*Meta)

	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return *res, nil
}

type Command struct {
	Cmd        string `json:"cmd"`
	DistrictID uint   `json:"districtId,omitempty"`
	PlaceID    uint   `json:"placeId,omitempty"`
	Batch      int    `json:"batch,omitempty"`
}

func EncodeCommands(cmds ...Command) string {
	data, err := json.Marshal(cmds)
	if err != nil {
		log.Print(e.Wrap("can't encode commands", err))
		return ""
	}
	return string(data)
}

func DecodeCommands(cmdText string) []Command {
	var cmds []Command
	err := json.Unmarshal([]byte(cmdText), &cmds)
	if err != nil {
		log.Print(e.Wrap("can't decode commands", err))
		return nil
	}
	return cmds
}
