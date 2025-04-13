package telegram

import (
	"SPBHistoryBot/clients/telegram"
	"SPBHistoryBot/lib/e"
	"log"
	"strconv"
	"strings"
)

const (
	HelpCmd         = "/help"
	StartCmd        = "/start"
	GetDistrictsCmd = "/getDistricts"
	DistrictCmd     = "/district"
	PlaceCmd        = "/place"
	DeleteCmd       = "/delete"
	CmdSplit        = `//`
	batchButSize    = 4
)

func (p *TgProcessor) doCmd(cmdText string, chatID int, username string) error {
	cmdText = strings.TrimSpace(cmdText)

	log.Printf("got new command: %s, from: %s", cmdText, username)
	switch cmdText {
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknown)
	}
}

func (p *TgProcessor) doCallbackCmd(cmdText string, chatID int, username string, messageID int) error {

	commands := splitCmds(cmdText)

	log.Printf("got new callback: %s, from: %s", cmdText, username)

	switch commands[0] {
	case StartCmd:
		return p.editToHello(chatID, messageID)

	case HelpCmd:
		return p.editToHelp(chatID, messageID)

	case GetDistrictsCmd:
		if len(commands) < 2 {
			return p.editToDistricts(chatID, messageID, 1)
		} else {
			batchNum, err := strconv.Atoi(commands[1])
			if err != nil {
				return e.Wrap("can't get batchNum", err)
			}
			return p.editToDistricts(chatID, messageID, batchNum)
		}

	case DistrictCmd:
		if len(commands) < 3 {
			districtId, err := strconv.Atoi(commands[1])
			if err != nil {
				return e.Wrap("can't get districtId", err)
			}
			return p.editToDistrict(chatID, messageID, districtId, 1)
		} else {
			batchNum, err := strconv.Atoi(commands[1])
			if err != nil {
				return e.Wrap("can't get batchNum", err)
			}
			districtId, err := strconv.Atoi(commands[2])
			if err != nil {
				return e.Wrap("can't get districtId", err)
			}
			return p.editToDistrict(chatID, messageID, districtId, batchNum)
		}
	default:
		return p.tg.SendMessage(chatID, callbackUnknown)
	}
}

func (p *TgProcessor) sendHello(chatID int) error {
	return p.tg.SendMessageWithButtons(chatID,
		msgHello,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: districtsBut, CallbackData: GetDistrictsCmd},
					{Text: hlpBut, CallbackData: HelpCmd},
				},
			},
		},
	)
}

func (p *TgProcessor) editToHelp(chatID int, messageID int) error {
	return p.tg.EditMessageWithButtons(chatID, messageID, msgHelp,
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
	return p.tg.EditMessageWithButtons(chatID, messageID,
		msgHello,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: districtsBut, CallbackData: GetDistrictsCmd},
					{Text: hlpBut, CallbackData: HelpCmd},
				},
			},
		},
	)
}

func (p *TgProcessor) editToDistricts(chatID int, messageID int, batchNum int) error {
	districts, err := p.storage.Districts()
	if err != nil {
		return e.Wrap("can't edit to districts", err)
	}
	kbs := make([][]telegram.InlineKeyboardButton, 0, batchButSize+1)

	for i := 0; i < batchNum*batchButSize && i < len(districts); i++ {
		kbs = append(kbs, []telegram.InlineKeyboardButton{{Text: districts[i].Name, CallbackData: joinCmds(DistrictCmd, strconv.FormatUint(uint64(districts[i].ID), 10))}})
	}
	if batchNum*batchButSize < len(districts) {
		kbs = append(kbs, []telegram.InlineKeyboardButton{{Text: nextBut, CallbackData: joinCmds(GetDistrictsCmd, strconv.Itoa(batchNum+1))}})
	}
	if batchNum > 1 {
		kbs = append(kbs, []telegram.InlineKeyboardButton{{Text: backBut, CallbackData: joinCmds(GetDistrictsCmd, strconv.Itoa(batchNum-1))}})
	} else {
		kbs = append(kbs, []telegram.InlineKeyboardButton{{Text: backBut, CallbackData: StartCmd}})
	}
	return p.tg.EditMessageWithButtons(chatID, messageID,
		msgChooseDistrict,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: kbs,
		},
	)
}

func (p *TgProcessor) editToDistrict(chatID int, messageID int, districtID int, batchNum int) error {
	district, err := p.storage.PickDistrict(districtID)
	if err != nil {
		return e.Wrap("can't edit to district", err)
	}
	keyboard := make([][]telegram.InlineKeyboardButton, 0, batchButSize+1)
	for i := 0; i < batchNum*batchButSize && i < len(district.Places); i++ {
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{{Text: district.Places[i].Name, CallbackData: joinCmds(DeleteCmd, PlaceCmd, district.Places[i].Name)}})
	}

	backForwardKeys := make([]telegram.InlineKeyboardButton, 0, 2)
	if batchNum*batchButSize < len(district.Places) {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: nextBut, CallbackData: joinCmds(GetDistrictsCmd, strconv.Itoa(batchNum+1))})
	}
	if batchNum > 1 {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: backBut, CallbackData: joinCmds(DistrictCmd, strconv.Itoa(batchNum-1))})
	} else {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: backBut, CallbackData: GetDistrictsCmd})
	}

	keyboard = append(keyboard, backForwardKeys)
	return p.tg.EditMessageWithButtons(chatID, messageID,
		msgChoosePlace,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	)
}

func joinCmds(commands ...string) string {
	var res string
	for i, cmd := range commands {
		res += cmd
		if i != len(commands)-1 {
			res += CmdSplit
		}
	}
	return res
}

func splitCmds(cmdText string) []string {
	return strings.Split(cmdText, CmdSplit)
}
