package telegram

import (
	"SPBHistoryBot/clients/telegram"
	"SPBHistoryBot/events"
	"SPBHistoryBot/lib/e"
	"SPBHistoryBot/lib/storage"
	"errors"
)

type Processor struct {
	tgFetcher telegram.UpdatesService
	tgSender  telegram.MessageService
	offset    int
	storage   storage.Storage
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
)

func NewProcessor(fetcher telegram.UpdatesService, sender telegram.MessageService, storage storage.Storage) *Processor {
	return &Processor{
		tgFetcher: fetcher,
		tgSender:  sender,
		storage:   storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tgFetcher.Updates(p.offset, limit)
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

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)

	case events.Callback:
		return p.processCallback(event)

	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := event.GetMeta()
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func (p *Processor) processCallback(event events.Event) error {
	meta, err := event.GetMeta()
	if err != nil {
		return e.Wrap("can't process message", err)
	}
	commands := splitCmds(event.Text)
	for i := 0; i < len(commands); i++ {
		if commands[i] != DeleteCmd { //За командой удаления сообщения могут идти другие команды
			if err := p.doCallbackCmd(joinCmds(commands[i:]...), meta.ChatId, meta.Username, meta.MessageID); err != nil {
				return e.Wrap("can't process message", err)
			}
			break

		} else {
			if err := p.doCallbackCmd(commands[i], meta.ChatId, meta.Username, meta.MessageID); err != nil {
				return e.Wrap("can't process message", err)
			}
		}

	}

	return nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message && upd.Message != nil {
		res.Meta = &events.Meta{
			ChatId:    upd.Message.Chat.ID,
			Username:  upd.Message.From.Username,
			MessageID: upd.Message.MessageID,
		}
	}
	if updType == events.Callback && upd.CallbackQuery != nil {
		res.Meta = &events.Meta{
			ChatId:    upd.CallbackQuery.From.ID,
			Username:  upd.CallbackQuery.From.Username,
			MessageID: upd.CallbackQuery.Message.MessageID,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	switch {
	case upd.Message != nil:
		return events.Message
	case upd.CallbackQuery != nil:
		return events.Callback
	default:
		return events.Unknown
	}
}

func fetchText(upd telegram.Update) string {
	if upd.Message != nil {
		return upd.Message.Text
	}
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}
	return ""
}
