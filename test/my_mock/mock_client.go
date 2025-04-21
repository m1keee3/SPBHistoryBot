package my_mock

import (
	"SPBHistoryBot/clients/telegram"
	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (m *Client) SendNoButtonsMessage(chatID int, text string) error {
	args := m.Called(chatID, text)
	return args.Error(0)
}

func (m *Client) SendMessage(chatID int, text string, markup telegram.InlineKeyboardMarkup) error {
	args := m.Called(chatID, text, markup)
	return args.Error(0)
}

func (m *Client) SendPhoto(chatID int, text string, photoURL string, markup telegram.InlineKeyboardMarkup) error {
	args := m.Called(chatID, text, photoURL, markup)
	return args.Error(0)
}

func (m *Client) EditMessage(chatID int, msgID int, text string, markup telegram.InlineKeyboardMarkup) error {
	args := m.Called(chatID, msgID, text, markup)
	return args.Error(0)
}

func (m *Client) EditPhoto(chatID int, msgID int, text string, photoURL string, markup telegram.InlineKeyboardMarkup) error {
	args := m.Called(chatID, msgID, text, photoURL, markup)
	return args.Error(0)
}

func (m *Client) DeleteMessage(chatID int, msgID int) error {
	args := m.Called(chatID, msgID)
	return args.Error(0)
}
