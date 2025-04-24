package test

import (
	tgClient "SPBHistoryBot/clients/telegram"
	"SPBHistoryBot/events"
	"SPBHistoryBot/events/telegram"
	"SPBHistoryBot/lib/storage"
	"SPBHistoryBot/test/my_mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessor_Process_StartCommand(t *testing.T) {
	client := new(my_mock.Client)
	st := new(my_mock.Storage)

	proc := telegram.NewProcessor(nil, client, st)

	event := events.Event{
		Type: events.Message,
		Text: "/start",
		Meta: &events.Meta{
			ChatID:   123,
			Username: "testuser",
		},
	}

	expectedKeyboard := tgClient.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{
			{
				{Text: telegram.LocationBut, CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.LocationHelpCmd})},
			},
			{
				{Text: telegram.DistrictsBut, CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.GetDistrictsCmd})},
				{Text: telegram.HelpBut, CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.HelpCmd})},
			},
		},
	}

	client.On("SendMessage", 123, mock.AnythingOfType("string"), expectedKeyboard).Return(nil)

	err := proc.Process(event)

	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestProcessor_Process_LocationHelpCommand(t *testing.T) {
	client := new(my_mock.Client)
	st := new(my_mock.Storage)

	proc := telegram.NewProcessor(nil, client, st)

	event := events.Event{
		Type: events.Callback,
		Text: events.EncodeCommands(events.Command{Cmd: telegram.LocationHelpCmd}),
		Meta: &events.Meta{
			ChatID:   123,
			Username: "testuser",
		},
	}

	expectedKeyboard := tgClient.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{
			{
				{Text: telegram.BackBut, CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.StartCmd})},
			},
		},
	}

	client.On("EditMessage", 123, mock.Anything, mock.AnythingOfType("string"), expectedKeyboard).Return(nil)

	err := proc.Process(event)

	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestProcessor_Process_GetDistrictsCommand(t *testing.T) {
	client := new(my_mock.Client)
	st := new(my_mock.Storage)

	proc := telegram.NewProcessor(nil, client, st)

	event := events.Event{
		Type: events.Callback,
		Text: events.EncodeCommands(events.Command{
			Cmd:   telegram.GetDistrictsCmd,
			Batch: 0,
		}),
		Meta: &events.Meta{
			ChatID:   123,
			Username: "testuser",
		},
	}

	districts := []storage.District{
		{Name: "District 1", ID: 1},
		{Name: "District 2", ID: 2},
	}

	st.On("Districts").Return(districts, nil)

	expectedKeyboard := tgClient.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{
			{{Text: "District 1", CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.GetPlacesCmd, DistrictID: 1})}},
			{{Text: "District 2", CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.GetPlacesCmd, DistrictID: 2})}},
			{
				{Text: telegram.BackBut, CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.StartCmd})},
			},
		},
	}

	client.On("EditMessage", 123, mock.Anything, mock.AnythingOfType("string"), expectedKeyboard).Return(nil)

	err := proc.Process(event)

	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestProcessor_Process_GetPlacesCommand(t *testing.T) {
	client := new(my_mock.Client)
	st := new(my_mock.Storage)

	proc := telegram.NewProcessor(nil, client, st)

	event := events.Event{
		Type: events.Callback,
		Text: events.EncodeCommands(events.Command{
			Cmd:        telegram.GetPlacesCmd,
			DistrictID: 1,
			Batch:      0,
		}),
		Meta: &events.Meta{
			ChatID:   123,
			Username: "testuser",
		},
	}

	district := storage.District{
		ID:   1,
		Name: "District 1",
		Places: []storage.Place{
			{Name: "Place 1", ID: 1},
			{Name: "Place 2", ID: 2},
		},
	}

	st.On("FindDistrict", 1).Return(&district, nil)

	expectedKeyboard := tgClient.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{
			{{Text: "Place 1", CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.DeleteCmd}, events.Command{Cmd: telegram.PlaceCmd, PlaceID: 1})}},
			{{Text: "Place 2", CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.DeleteCmd}, events.Command{Cmd: telegram.PlaceCmd, PlaceID: 2})}},
			{
				{Text: telegram.BackBut, CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.GetDistrictsCmd})},
			},
		},
	}

	client.On("EditMessage", 123, mock.Anything, mock.AnythingOfType("string"), expectedKeyboard).Return(nil)

	err := proc.Process(event)

	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestProcessor_Process_SendDistrictCommand(t *testing.T) {
	client := new(my_mock.Client)
	st := new(my_mock.Storage)

	proc := telegram.NewProcessor(nil, client, st)

	event := events.Event{
		Type: events.Callback,
		Text: events.EncodeCommands(events.Command{
			Cmd:        telegram.SendDistrictCmd,
			DistrictID: 1,
		}),
		Meta: &events.Meta{
			ChatID:   123,
			Username: "testuser",
		},
	}

	district := storage.District{
		ID:   1,
		Name: "District 1",
		Places: []storage.Place{
			{Name: "Place 1", ID: 1},
		},
	}

	st.On("FindDistrict", 1).Return(&district, nil)

	expectedKeyboard := tgClient.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{
			{
				{Text: "Place 1", CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.DeleteCmd}, events.Command{Cmd: telegram.PlaceCmd, PlaceID: 1})},
			},
			{
				{Text: telegram.BackBut, CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.GetDistrictsCmd})},
			},
		},
	}

	client.On("SendMessage", 123, mock.AnythingOfType("string"), expectedKeyboard).Return(nil)

	err := proc.Process(event)

	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestProcessor_Process_PlaceCommand(t *testing.T) {
	client := new(my_mock.Client)
	st := new(my_mock.Storage)

	proc := telegram.NewProcessor(nil, client, st)

	event := events.Event{
		Type: events.Callback,
		Text: events.EncodeCommands(events.Command{
			Cmd:     telegram.PlaceCmd,
			PlaceID: 1,
		}),
		Meta: &events.Meta{
			ChatID:   123,
			Username: "testuser",
		},
	}

	place := storage.Place{
		ID:         1,
		Name:       "Place 1",
		Text:       "Description of Place 1",
		Image:      "http://example.com/image.jpg",
		DistrictID: 1,
	}

	st.On("FindPlace", 1).Return(&place, nil)

	expectedKeyboard := tgClient.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{
			{
				{Text: telegram.BackBut, CallbackData: events.EncodeCommands(events.Command{Cmd: telegram.DeleteCmd}, events.Command{Cmd: telegram.SendDistrictCmd, DistrictID: place.DistrictID})},
			},
		},
	}

	client.On("SendPhoto", 123, mock.AnythingOfType("string"), place.Image, expectedKeyboard).Return(nil)

	err := proc.Process(event)

	require.NoError(t, err)
	client.AssertExpectations(t)
}
