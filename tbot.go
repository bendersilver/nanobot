package nanobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Bad int

const (
	OK Bad = iota
	BadChat
	BadToken
	BadOther
)
const tURL = "https://api.telegram.org/bot%s/"

type Result struct {
	Desc   string
	Status Bad
	Code   int
	ID     int64
}

type response struct {
	OK   bool   `json:"ok"`
	Code int    `json:"error_code"`
	Desc string `json:"description"`
	Res  struct {
		ID int64 `json:"message_id"`
	} `json:"result"`
}

type Body struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
	Mode   string `json:"parse_mode,omitempty"`
}

type Bot struct {
	Token, uri string
}

func New(token string) (*Bot, error) {
	var b Bot
	b.Token = token
	b.uri = fmt.Sprintf(tURL, token)
	r := b.getMe()
	if r.Status != OK {
		return nil, fmt.Errorf(r.Desc)
	}
	return &b, nil
}

func (b *Bot) SendMessage(arg *Body) *Result {
	r := b.req("sendMessage", arg)

	if r.Status != OK {
		switch r.Desc {
		case "Bad Request: chat not found":
			r.Status = BadChat
		}

	}
	return r
}

func (b *Bot) getMe() *Result {
	r := b.req("getMe", nil)
	return r
}

func (b *Bot) req(met string, arg *Body) (res *Result) {
	res = new(Result)
	var resp *http.Response
	var err error
	defer func() {
		if err != nil {
			res.Status = BadOther
			res.Desc = err.Error()
		}
	}()

	if arg == nil {
		resp, err = http.Get(b.uri + met)

	} else {
		var r bytes.Buffer
		err = json.NewEncoder(&r).Encode(&arg)
		if err != nil {
			return
		}

		resp, err = http.Post(b.uri+met, "application/json", &r)

	}
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var body response
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return
	}

	if !body.OK {
		res.Desc = body.Desc
		res.Code = body.Code
		res.Status = BadOther
		if body.Code == 401 || body.Code == 403 {
			res.Status = BadToken
		}
	} else {
		res.Status = OK
		res.ID = body.Res.ID
	}
	return
}
