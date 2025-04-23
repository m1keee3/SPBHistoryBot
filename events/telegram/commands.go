package telegram

import (
	"SPBHistoryBot/clients/telegram"
	"SPBHistoryBot/events"
	"SPBHistoryBot/lib/e"
	"log"
)

const (
	HelpCmd            = "/help"
	StartCmd           = "/start"
	GetDistrictsCmd    = "/getDistricts"
	GetPlacesCmd       = "/getPlaces"
	SendDistrictCmd    = "/sendDistrict"
	PlaceCmd           = "/place"
	RequestLocationCmd = "/requestLocation"
	LocationHelpCmd    = "/locationHelp"
	DeleteCmd          = "/delete"
	batchButSize       = 4
)

func (p *Processor) doCmd(cmd string, meta events.Meta) error {

	log.Printf("got new command: %s, from: %s", cmd, meta.Username)

	switch cmd {
	case StartCmd:
		return p.sendHello(meta.ChatID)

	case RequestLocationCmd:
		nearPlace, err := p.storage.FindNearPlace(meta.Latitude, meta.Longitude)
		if err != nil {
			return err
		}
		return p.sendPlace(meta.ChatID, int(nearPlace.ID))

	default:
		return p.tgSender.SendNoButtonsMessage(meta.ChatID, MsgUnknown)
	}
}

func (p *Processor) doCallbackCmd(cmd events.Command, meta events.Meta) error {

	log.Printf("got new callback: %s, from: %s", cmd.Cmd, meta.Username)

	switch cmd.Cmd {
	case StartCmd:
		return p.editToHello(meta.ChatID, meta.MessageID)

	case LocationHelpCmd:
		return p.editToLocationHelp(meta.ChatID, meta.MessageID)

	case HelpCmd:
		return p.editToHelp(meta.ChatID, meta.MessageID)

	case GetDistrictsCmd:
		return p.editToDistricts(meta.ChatID, meta.MessageID, cmd.Batch)

	case GetPlacesCmd:
		return p.editToDistrict(meta.ChatID, meta.MessageID, int(cmd.DistrictID), cmd.Batch)

	case SendDistrictCmd:
		return p.sendDistrict(meta.ChatID, int(cmd.DistrictID))

	case PlaceCmd:
		return p.sendPlace(meta.ChatID, int(cmd.PlaceID))

	case DeleteCmd:
		return p.tgSender.DeleteMessage(meta.ChatID, meta.MessageID)

	default:
		return p.tgSender.SendNoButtonsMessage(meta.ChatID, CallbackUnknown)
	}
}

func (p *Processor) sendHello(chatID int) error {

	locationHelpCmd := events.Command{Cmd: LocationHelpCmd}
	getDistrictsCmd := events.Command{Cmd: GetDistrictsCmd}
	helpCmd := events.Command{Cmd: HelpCmd}

	return p.tgSender.SendMessage(chatID,
		MsgHello,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: LocationBut, CallbackData: events.EncodeCommands(locationHelpCmd)},
				},
				{
					{Text: DistrictsBut, CallbackData: events.EncodeCommands(getDistrictsCmd)},
					{Text: HelpBut, CallbackData: events.EncodeCommands(helpCmd)},
				},
			},
		},
	)
}

func (p *Processor) editToLocationHelp(chatID int, messageID int) error {

	startCmd := events.Command{Cmd: StartCmd}

	return p.tgSender.EditMessage(chatID, messageID, MsgLocationHelp,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: BackBut, CallbackData: events.EncodeCommands(startCmd)},
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

	locationHelpCmd := events.Command{Cmd: LocationHelpCmd}
	getDistrictsCmd := events.Command{Cmd: GetDistrictsCmd}
	helpCmd := events.Command{Cmd: HelpCmd}

	return p.tgSender.EditMessage(chatID, messageID,
		MsgHello,
		telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: LocationBut, CallbackData: events.EncodeCommands(locationHelpCmd)},
				},
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
