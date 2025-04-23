package test

import (
	tgClient "SPBHistoryBot/clients/telegram"
	"SPBHistoryBot/events"
	"SPBHistoryBot/events/telegram"
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
