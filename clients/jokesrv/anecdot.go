package jokesrv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tg_ics_useful_bot/lib/e"
)

type Joke struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

var jokesrvURL = "https://jokesrv.rubedo.cloud/"

func Anecdot() (string, error) {
	c := http.Client{}
	req, err := http.NewRequest(http.MethodGet, jokesrvURL, nil)
	if err != nil {
		return "", e.Wrap(fmt.Sprintf("[ERROR] can't make request, url %s: ", jokesrvURL), err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return "", e.Wrap("[ERROR] can't do request: ", err)
	}
	defer func() { _ = resp.Body.Close() }()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", e.Wrap("[ERROR] can't get response: ", err)
	}
	var anecdot Joke
	err = json.Unmarshal(data, &anecdot)
	return anecdot.Content, nil
}
