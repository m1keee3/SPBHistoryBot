package telegram

import (
	"log"
	"strings"
)

const (
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *TgProcessor) doCmd(text string, chatID int, username string) error {
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

func (p *TgProcessor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *TgProcessor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}
