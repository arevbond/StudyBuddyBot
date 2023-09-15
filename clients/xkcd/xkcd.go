package xkcd

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"tg_ics_useful_bot/lib/e"
	"time"
)

const (
	baseHost = "xkcd.com"
	baseEnd  = "/info.0.json"
)

const (
	minComicsNumber       = 1
	NotExistsComicsNumber = 404
	maxComicsNumber       = 2825
)

type Client struct {
	Host    string
	BaseEnd string
	Client  http.Client
}

func New() *Client {
	return &Client{
		Host:    baseHost,
		BaseEnd: baseEnd,
		Client:  http.Client{},
	}
}

func RandomComics() (Comics, error) {
	c := New()
	number := randomComicsNumber()
	comics, err := c.comics(number)
	if err != nil {
		return Comics{}, e.Wrap("can't get comics: ", err)
	}
	return comics, nil
}

func randomComicsNumber() int {
	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(maxComicsNumber)
	if number == 0 || number == NotExistsComicsNumber {
		number++
	}
	return number
}

func (c *Client) comics(num int) (Comics, error) {
	data, err := c.doRequest(num)
	if err != nil {
		return Comics{}, e.Wrap("can't get comics: ", err)
	}

	var comics Comics

	err = json.Unmarshal(data, &comics)
	if err != nil {
		return Comics{}, e.Wrap("can't unmarshall json to Comics: ", err)
	}
	return comics, nil
}

func (c *Client) doRequest(num int) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   baseHost,
		Path:   strconv.Itoa(num) + baseEnd,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap("can't do request: ", err)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, e.Wrap("can't do request: ", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't do request: ", err)
	}
	return body, nil
}
