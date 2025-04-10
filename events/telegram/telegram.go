package telegram

import (
	"SPBHistoryBot/clients/telegram"
	"SPBHistoryBot/events"
	"SPBHistoryBot/lib/e"
	"errors"
)

type TgProcessor struct {
	tg     *telegram.Client
	offset int
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
)

func NewProcessor(client *telegram.Client) *TgProcessor {
	return &TgProcessor{
		tg: client,
	}
}

func (p *TgProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("failed to fetch telegram updates", err)
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

func (p *TgProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *TgProcessor) processMessage(event events.Event) error {
	meta, err := event.GetMeta()
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = &events.Meta{
			ChatId:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
