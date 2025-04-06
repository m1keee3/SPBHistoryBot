package telegram

import (
	"SPBHistoryBot/clients"
	"SPBHistoryBot/lib/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod = "getUpdates"
)

type TelegramClient struct {
	host     string
	basePath string
	client   *http.Client
}

func NewTelegramClient(host string, basePath string) TelegramClient {
	return TelegramClient{
		host:     host,
		basePath: newBasePath(basePath),
		client:   &http.Client{},
	}
}

func newBasePath(basePath string) string {
	return "bot" + basePath
}

func (c *TelegramClient) Updates(offset int, limit int) ([]clients.Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res clients.UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *TelegramClient) doRequest(method string, query url.Values) (data []byte, err error) {
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

func (c *TelegramClient) SendMessage() {

}
