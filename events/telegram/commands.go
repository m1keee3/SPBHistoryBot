package telegram

import (
	"SPBHistoryBot/clients/telegram"
	"SPBHistoryBot/events"
	"SPBHistoryBot/lib/e"
	"log"
)

const (
	HelpCmd         = "/help"
	StartCmd        = "/start"
	GetDistrictsCmd = "/getDistricts"
	GetPlacesCmd    = "/getPlaces"
	SendDistrictCmd = "/sendDistrict"
	PlaceCmd        = "/place"
	DeleteCmd       = "/delete"
	batchButSize    = 4
)

func (p *Processor) doCmd(cmd string, chatID int, username string) error {

	log.Printf("got new command: %s, from: %s", cmd, username)

	switch cmd {
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tgSender.SendNoButtonsMessage(chatID, MsgUnknown)
	}
}

func (p *Processor) doCallbackCmd(cmd events.Command, chatID int, username string, messageID int) error {

	log.Printf("got new callback: %s, from: %s", cmd.Cmd, username)

	switch cmd.Cmd {
	case StartCmd:
		return p.editToHello(chatID, messageID)

	case HelpCmd:
		return p.editToHelp(chatID, messageID)

	case GetDistrictsCmd:
		return p.editToDistricts(chatID, messageID, cmd.Batch)

	case GetPlacesCmd:
		return p.editToDistrict(chatID, messageID, int(cmd.DistrictID), cmd.Batch)

	case SendDistrictCmd:
		return p.sendDistrict(chatID, int(cmd.DistrictID))

	case PlaceCmd:
		return p.sendPlace(chatID, int(cmd.PlaceID))

	case DeleteCmd:
		return p.tgSender.DeleteMessage(chatID, messageID)

	default:
		return p.tgSender.SendNoButtonsMessage(chatID, CallbackUnknown)
	}
}

func (p *Processor) sendHello(chatID int) error {

	getDistrictsCmd := events.Command{Cmd: GetDistrictsCmd}
	helpCmd := events.Command{Cmd: HelpCmd}

	return p.tgSender.SendMessage(chatID,
		MsgHello,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: DistrictsBut, CallbackData: events.EncodeCommands(getDistrictsCmd)},
					{Text: HelpBut, CallbackData: events.EncodeCommands(helpCmd)},
				},
			},
		},
	)
}

func (p *Processor) editToHelp(chatID int, messageID int) error {

	startCmd := events.Command{Cmd: StartCmd}

	return p.tgSender.EditMessage(chatID, messageID, MsgHelp,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: BackBut, CallbackData: events.EncodeCommands(startCmd)},
				},
			},
		},
	)
}

func (p *Processor) editToHello(chatID int, messageID int) error {

	getDistrictsCmd := events.Command{Cmd: GetDistrictsCmd}
	helpCmd := events.Command{Cmd: HelpCmd}

	return p.tgSender.EditMessage(chatID, messageID,
		MsgHello,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: DistrictsBut, CallbackData: events.EncodeCommands(getDistrictsCmd)},
					{Text: HelpBut, CallbackData: events.EncodeCommands(helpCmd)},
				},
			},
		},
	)
}

func (p *Processor) editToDistricts(chatID int, messageID int, batchNum int) error {
	districts, err := p.storage.Districts()
	if err != nil {
		return e.Wrap("can't edit to districts", err)
	}

	startCmd := events.Command{Cmd: StartCmd}
	keyboard := make([][]telegram.InlineKeyboardButton, 0, batchButSize+1)

	for i := batchNum * batchButSize; i < (batchNum+1)*batchButSize && i < len(districts); i++ {
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{{Text: districts[i].Name, CallbackData: events.EncodeCommands(events.Command{
			Cmd:        GetPlacesCmd,
			DistrictID: districts[i].ID,
		})}})
	}

	backForwardKeys := make([]telegram.InlineKeyboardButton, 0, 2)

	if batchNum > 0 {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: BackBut, CallbackData: events.EncodeCommands(events.Command{
			Cmd:   GetDistrictsCmd,
			Batch: batchNum - 1,
		})})
	} else {
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{{Text: BackBut, CallbackData: events.EncodeCommands(startCmd)}})
	}

	if (batchNum+1)*batchButSize < len(districts) {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: NextBut, CallbackData: events.EncodeCommands(events.Command{
			Cmd:   GetDistrictsCmd,
			Batch: batchNum + 1,
		})})
	}

	keyboard = append(keyboard, backForwardKeys)
	return p.tgSender.EditMessage(chatID, messageID,
		MsgChooseDistrict,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	)
}

func (p *Processor) editToDistrict(chatID int, messageID int, districtID int, batchNum int) error {
	district, err := p.storage.FindDistrict(districtID)
	if err != nil {
		return e.Wrap("can't edit to district", err)
	}

	deleteCmd := events.Command{Cmd: DeleteCmd}
	keyboard := make([][]telegram.InlineKeyboardButton, 0, batchButSize+1)

	for i := batchNum * batchButSize; i < (batchNum+1)*batchButSize && i < len(district.Places); i++ {
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{{Text: district.Places[i].Name, CallbackData: events.EncodeCommands(deleteCmd, events.Command{
			Cmd:     PlaceCmd,
			PlaceID: district.Places[i].ID,
		})}})
	}

	backForwardKeys := make([]telegram.InlineKeyboardButton, 0, 2)

	if batchNum > 0 {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: BackBut, CallbackData: events.EncodeCommands(events.Command{
			Cmd:        GetPlacesCmd,
			DistrictID: uint(districtID),
			Batch:      batchNum - 1,
		})})
	} else {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: BackBut, CallbackData: events.EncodeCommands(events.Command{
			Cmd: GetDistrictsCmd,
		})})
	}

	if (batchNum+1)*batchButSize < len(district.Places) {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: NextBut, CallbackData: events.EncodeCommands(events.Command{
			Cmd:        GetPlacesCmd,
			DistrictID: uint(districtID),
			Batch:      batchNum + 1,
		})})
	}

	keyboard = append(keyboard, backForwardKeys)
	return p.tgSender.EditMessage(chatID, messageID,
		MsgChoosePlace,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	)
}

func (p *Processor) sendDistrict(chatID int, districtID int) error {
	district, err := p.storage.FindDistrict(districtID)
	if err != nil {
		return e.Wrap("can't edit to district", err)
	}

	deleteCmd := events.Command{Cmd: DeleteCmd}
	keyboard := make([][]telegram.InlineKeyboardButton, 0, batchButSize+1)

	for i := 0; i < batchButSize && i < len(district.Places); i++ {
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{{Text: district.Places[i].Name, CallbackData: events.EncodeCommands(deleteCmd, events.Command{
			Cmd:     PlaceCmd,
			PlaceID: district.Places[i].ID,
		})}})
	}

	backForwardKeys := make([]telegram.InlineKeyboardButton, 0, 2)
	backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: BackBut, CallbackData: events.EncodeCommands(events.Command{Cmd: GetDistrictsCmd})})

	if batchButSize < len(district.Places) {
		backForwardKeys = append(backForwardKeys, telegram.InlineKeyboardButton{Text: NextBut, CallbackData: events.EncodeCommands(events.Command{
			Cmd:        GetPlacesCmd,
			DistrictID: uint(districtID),
			Batch:      1,
		})})
	}

	keyboard = append(keyboard, backForwardKeys)
	return p.tgSender.SendMessage(chatID,
		MsgChoosePlace,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	)
}

func (p *Processor) sendPlace(chatID int, placeID int) error {
	place, err := p.storage.FindPlace(placeID)
	if err != nil {
		return e.Wrap("can't find a place", err)
	}

	deleteCmd := events.Command{Cmd: DeleteCmd}
	sendDistrictCmd := events.Command{
		Cmd:        SendDistrictCmd,
		DistrictID: place.DistrictID,
	}

	return p.tgSender.SendPhoto(chatID,
		place.Name+"\n\n"+place.Text,
		place.Image,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: BackBut, CallbackData: events.EncodeCommands(deleteCmd, sendDistrictCmd)},
				},
			},
		},
	)
}
