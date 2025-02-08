package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log-proj/pkg/lib/e"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const getUpdatesMethod = "getUpdates"
const sendMessageEndpoint = "sendMessage"

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Updates, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q, nil)
	if err != nil {
		return nil, e.Warp("can`t get updates", err)
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Warp("can`t unmarshal updates", err)
	}

	return res.Result, nil
}

type keyboard interface {
}

func (c *Client) SendMessage(chatID int, text string, keyboard keyboard) error {
	msg := OutcomingMessage{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: keyboard,
	}

	// Преобразуем сообщение в JSON
	body, err := json.Marshal(msg)
	if err != nil {
		return e.Warp("can't marshal message", err)
	}

	// Отправляем запрос
	_, err = c.doRequest(sendMessageEndpoint, nil, body)
	if err != nil {
		return e.Warp("can't send message", err)
	}

	return nil
}

func (c *Client) doRequest(Endpoint string, q url.Values, body []byte) ([]byte, error) {
	const errMsg = "can`t do request"
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, Endpoint),
	}
	method := http.MethodGet
	var buf io.Reader

	if body != nil {
		buf = bytes.NewBuffer(body)
		method = http.MethodPost
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, e.Warp(errMsg, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Warp(errMsg, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, e.Warp(errMsg, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, b))
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Warp(errMsg, err)
	}

	return body, nil
}
