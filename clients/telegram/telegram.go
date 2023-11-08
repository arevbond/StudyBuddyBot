package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"tg_ics_useful_bot/lib/e"
	"time"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const timeToBan = 120

const (
	getUpdatesMethod            = "getUpdates"
	sendMessageMethod           = "sendMessage"
	sendPhotoMethod             = "sendPhoto"
	deleteMessageMethod         = "deleteMessage"
	banChatMemberMethod         = "banChatMember"
	getChatAdministratorsMethod = "getChatAdministrators"
)

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

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't get updates: ", err) }()

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

func (c *Client) ChatAdministrators(chatID int) ([]User, error) {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))

	data, err := c.doRequest(getChatAdministratorsMethod, q)
	if err != nil {
		return nil, e.Wrap("can't get chat administrators: ", err)
	}
	var dataResponse ChatMemberAdministratorResponse

	if err := json.Unmarshal(data, &dataResponse); err != nil {
		return nil, err
	}

	var result []User
	for _, admin := range dataResponse.Result {
		result = append(result, admin.User)
	}

	return result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("parse_mode", "Markdown")

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) SendPhoto(chatID int, urlPhoto string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("photo", urlPhoto)

	_, err := c.doRequest(sendPhotoMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) DeleteMessage(chatID int, messageID int) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(messageID))

	_, err := c.doRequest(deleteMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) BanChatMember(chatID int, userID int, timeout int) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("user_id", strconv.Itoa(userID))
	q.Add("until_date", strconv.Itoa(int(time.Now().Unix())+timeToBan))

	_, err := c.doRequest(banChatMemberMethod, q)
	if err != nil {
		return e.Wrap("can't ban user: ", err)
	}

	return nil
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
