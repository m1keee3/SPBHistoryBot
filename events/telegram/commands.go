package telegram

import (
	"SPBHistoryBot/clients/telegram"
	"log"
	"strings"
)

const (
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *TgProcessor) sendMsg(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command: %s, from: %s", text, username)
	switch text {
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknown)
	}
}

func (p *TgProcessor) sendCallback(text string, chatID int, username string, messageID int) error {
	text = strings.TrimSpace(text)

	log.Printf("got new callback: %s, from: %s", text, username)

	switch text {
	case HelpCmd:
		return p.editToHelp(chatID, messageID)
	case StartCmd:
		return p.editToHello(chatID, messageID)
	default:
		return p.tg.SendMessage(chatID, msgUnknown)
	}
}

func (p *TgProcessor) sendHello(chatID int) error {
	return p.tg.SendMessageWithButtons(chatID,
		msgHello,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: hlpBut, CallbackData: HelpCmd},
				},
			},
		},
	)
}

func (p *TgProcessor) sendHelp(chatID int) error {
	return p.tg.SendMessageWithButtons(chatID,
		msgHelp,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: backBut, CallbackData: StartCmd},
				},
			},
		},
	)
}

func (p *TgProcessor) editToHelp(chatID int, messageID int) error {
	return p.tg.EditMessageTextWithButtons(chatID, messageID, msgHelp,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: backBut, CallbackData: StartCmd},
				},
			},
		},
	)
}

func (p *TgProcessor) editToHello(chatID int, messageID int) error {
	return p.tg.EditMessageTextWithButtons(chatID, messageID, msgHello,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: hlpBut, CallbackData: HelpCmd},
				},
			},
		},
	)
}
