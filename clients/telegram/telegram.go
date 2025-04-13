package telegram

import (
	"SPBHistoryBot/lib/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod    = "getUpdates"
	sendMessageMethod   = "sendMessage"
	editMessageMethod   = "editMessageText"
	sendPhotoMethod     = "sendPhoto"
	editPhotoMethod     = "editMessageMedia"
	deleteMessageMethod = "deleteMessage"
)

type Client struct {
	host     string
	basePath string
	client   *http.Client
}

func NewClient(host string, basePath string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(basePath),
		client:   &http.Client{},
	}
}

func newBasePath(basePath string) string {
	return "bot" + basePath
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) SendMessageWithButtons(chatID int, text string, keyboard InlineKeyboardMarkup) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	kbjson, err := json.Marshal(keyboard)
	if err != nil {
		return e.Wrap("can't marshal keyboard to json", err)
	}
	q.Add("reply_markup", string(kbjson))

	if _, err := c.doRequest(sendMessageMethod, q); err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) SendPhotoWithButtons(chatID int, text string, photoURL string, keyboard InlineKeyboardMarkup) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("caption", text)
	q.Add("photo", photoURL)

	kbjson, err := json.Marshal(keyboard)
	if err != nil {
		return e.Wrap("can't marshal keyboard to json", err)
	}
	q.Add("reply_markup", string(kbjson))

	if _, err := c.doRequest(sendPhotoMethod, q); err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) EditMessageWithButtons(chatID int, messageID int, text string, markup InlineKeyboardMarkup) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(messageID))
	q.Add("text", text)

	data, err := json.Marshal(markup)
	if err != nil {
		return e.Wrap("can't marshal inline keyboard", err)
	}

	q.Add("reply_markup", string(data))

	_, err = c.doRequest(editMessageMethod, q)
	if err != nil {
		return e.Wrap("can't edit message with buttons", err)
	}

	return nil
}

func (c *Client) EditPhotoWithButtons(chatID int, messageID int, text string, photoURL string, markup InlineKeyboardMarkup) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(messageID))
	q.Add("caption", text)
	q.Add("photo", photoURL)

	data, err := json.Marshal(markup)
	if err != nil {
		return e.Wrap("can't marshal inline keyboard", err)
	}

	q.Add("reply_markup", string(data))

	if _, err := c.doRequest(editPhotoMethod, q); err != nil {
		return e.Wrap("can't edit message with buttons", err)
	}

	return nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	if _, err := c.doRequest(sendMessageMethod, q); err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) DeleteMessage(chatID int, messageID int) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(messageID))

	if _, err := c.doRequest(deleteMessageMethod, q); err != nil {
		return e.Wrap("can't delete message", err)
	}

	return nil
}
