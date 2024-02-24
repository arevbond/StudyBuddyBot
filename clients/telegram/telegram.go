package telegram

import (
	"bytes"
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
	AdminsID []int
}

type ParseMode string

const (
	WithoutParseMode ParseMode = ""
	MarkdownV2       ParseMode = "MarkdownV2"
	HTML             ParseMode = "HTML"
	Markdown         ParseMode = "Markdown"
)

const timeToBan = 120

const (
	getUpdatesMethod            = "getUpdates"
	sendMessageMethod           = "sendMessage"
	sendPhotoMethod             = "sendPhoto"
	deleteMessageMethod         = "deleteMessage"
	banChatMemberMethod         = "banChatMember"
	getChatAdministratorsMethod = "getChatAdministrators"
	sendPoll                    = "sendPoll"
)

func New(host string, token string, adminsID []int) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
		AdminsID: adminsID,
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
	data, err := c.doRequestWithQuery(getUpdatesMethod, q)
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

	data, err := c.doRequestWithQuery(getChatAdministratorsMethod, q)
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

func (c *Client) SendMessage(chatID int, text string, parseMode ParseMode, replyToMessageID int) error {
	message := Message{chatID, text, string(parseMode), replyToMessageID}
	jsonData, err := json.Marshal(message)
	if err != nil {
		return e.Wrap("can't convert message to json: ", err)
	}
	_, err = c.doRequestWithBody(sendMessageMethod, jsonData)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) SendPhoto(chatID int, urlPhoto string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("photo", urlPhoto)

	_, err := c.doRequestWithQuery(sendPhotoMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) DeleteMessage(chatID int, messageID int) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(messageID))

	_, err := c.doRequestWithQuery(deleteMessageMethod, q)
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

	_, err := c.doRequestWithQuery(banChatMemberMethod, q)
	if err != nil {
		return e.Wrap("can't ban user: ", err)
	}

	return nil
}

func (c *Client) doRequestWithQuery(method string, query url.Values) (data []byte, err error) {
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

func (c *Client) doRequestWithBody(method string, message []byte) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request with json", err) }()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(message))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	resultBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return resultBody, nil
}

func (c *Client) SendPoll(poll SendPoll) error {
	jsonData, err := json.Marshal(poll)
	if err != nil {
		return e.Wrap("can't convert poll to json: ", err)
	}
	_, err = c.doRequestWithBody(sendPoll, jsonData)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}
