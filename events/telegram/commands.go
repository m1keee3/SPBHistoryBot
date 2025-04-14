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
	SendDistrictCmd = "/sendDistrict"
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
		return p.tg.SendNoButtonsMessage(chatID, msgUnknown)
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
		}

		batchNum, err := strconv.Atoi(commands[1])
		if err != nil {
			return e.Wrap("can't get batchNum", err)
		}

		return p.editToDistricts(chatID, messageID, batchNum)

	case DistrictCmd:
		if len(commands) < 2 {
			return p.tg.SendNoButtonsMessage(chatID, "Некорректная команда района")
		}
		districtId, err := strconv.Atoi(commands[1])
		if err != nil {
			return e.Wrap("can't get districtId", err)
		}
		batchNum := 1
		if len(commands) > 2 {
			batchNum, err = strconv.Atoi(commands[2])
			if err != nil {
				return e.Wrap("can't get batchNum", err)
			}
		}
		return p.editToDistrict(chatID, messageID, districtId, batchNum)

	case SendDistrictCmd:
		districtId, err := strconv.Atoi(commands[1])
		if err != nil {
			return e.Wrap("can't get districtId", err)
		}
		return p.sendDistrict(chatID, districtId)

	case PlaceCmd:
		placeID, err := strconv.Atoi(commands[1])
		if err != nil {
			return e.Wrap("can't get placeId", err)
		}
		return p.sendPlace(chatID, placeID)

	case DeleteCmd:
		return p.tg.DeleteMessage(chatID, messageID)

	default:
		return p.tg.SendNoButtonsMessage(chatID, callbackUnknown)
	}
}

func (p *TgProcessor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID,
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
	return p.tg.EditMessage(chatID, messageID, msgHelp,
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
	return p.tg.EditMessage(chatID, messageID,
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

	for i := (batchNum - 1) * batchButSize; i < batchNum*batchButSize && i < len(districts); i++ {
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
	return p.tg.EditMessage(chatID, messageID,
		msgChooseDistrict,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: kbs,
		},
	)
}

func (p *TgProcessor) editToDistrict(chatID int, messageID int, districtID int, batchNum int) error {
	district, err := p.storage.FindDistrict(districtID)
	if err != nil {
		return e.Wrap("can't edit to district", err)
	}
	keyboard := make([][]telegram.InlineKeyboardButton, 0, batchButSize+1)
	for i := (batchNum - 1) * batchButSize; i < batchNum*batchButSize && i < len(district.Places); i++ {
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{{Text: district.Places[i].Name, CallbackData: joinCmds(DeleteCmd, PlaceCmd, strconv.FormatUint(uint64(district.Places[i].ID), 10))}})
	}

	backForwardKeys := make([]telegram.InlineKeyboardButton, 0, 2)
	if batchNum > 1 {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: backBut, CallbackData: joinCmds(DistrictCmd, strconv.FormatUint(uint64(districtID), 10), strconv.Itoa(batchNum-1))})
	} else {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: backBut, CallbackData: GetDistrictsCmd})
	}

	if batchNum*batchButSize < len(district.Places) {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: nextBut, CallbackData: joinCmds(DistrictCmd, strconv.FormatUint(uint64(districtID), 10), strconv.Itoa(batchNum+1))})
	}

	keyboard = append(keyboard, backForwardKeys)
	return p.tg.EditMessage(chatID, messageID,
		msgChoosePlace,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	)
}

func (p *TgProcessor) sendDistrict(chatID int, districtID int) error {
	district, err := p.storage.FindDistrict(districtID)
	if err != nil {
		return e.Wrap("can't edit to district", err)
	}
	keyboard := make([][]telegram.InlineKeyboardButton, 0, batchButSize+1)
	for i := 0; i < batchButSize && i < len(district.Places); i++ {
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{{Text: district.Places[i].Name, CallbackData: joinCmds(DeleteCmd, PlaceCmd, strconv.FormatUint(uint64(district.Places[i].ID), 10))}})
	}

	backForwardKeys := make([]telegram.InlineKeyboardButton, 0, 2)
	backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: backBut, CallbackData: GetDistrictsCmd})

	if batchButSize < len(district.Places) {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: nextBut, CallbackData: joinCmds(DistrictCmd, strconv.FormatUint(uint64(districtID), 10), strconv.Itoa(2))})
	}

	keyboard = append(keyboard, backForwardKeys)
	return p.tg.SendMessage(chatID,
		msgChoosePlace,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	)
}

func (p *TgProcessor) sendPlace(chatID int, placeID int) error {
	place, err := p.storage.FindPlace(placeID)
	if err != nil {
		return e.Wrap("can't find a place", err)
	}

	return p.tg.SendPhoto(chatID,
		place.Name+"\n\n"+place.Text,
		place.Image,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: backBut, CallbackData: joinCmds(DeleteCmd, SendDistrictCmd, strconv.FormatUint(uint64(place.DistrictID), 10))},
				},
			},
		},
	)
}

func joinCmds(commands ...string) string {
	return strings.Join(commands, CmdSplit)
}

func splitCmds(cmdText string) []string {
	return strings.Split(cmdText, CmdSplit)
}
