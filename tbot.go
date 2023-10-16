package nanobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Bad -
type Bad int

const (
	// OK -
	OK Bad = iota
	// BadChat -
	BadChat
	// BadToken -
	BadToken
	// BadToken -
	Sleep
	// BadOther -
	BadOther
)
const tURL = "https://api.telegram.org/bot%s/"

// Result -
type Result struct {
	Desc   string
	Status Bad
	Code   int
	ID     int64
}

type response struct {
	OK    bool            `json:"ok"`
	Code  int             `json:"error_code"`
	Desc  string          `json:"description"`
	Res   json.RawMessage `json:"result"`
	Parms struct {
		MtoChatID  int64 `json:"migrate_to_chat_id"`
		RetryAfter int   `json:"retry_after"`
	} `json:"parameters"`
}

// Body -
type Body struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
	Mode   string `json:"parse_mode,omitempty"`
}

// Bot -
type Bot struct {
	Token, uri string
}

// New -
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

// DeleteMessage -
func (b *Bot) DeleteMessage(chatID, messageID int64) *Result {
	r, _ := b.req("deleteMessage", map[string]int64{
		"chat_id":    chatID,
		"message_id": messageID,
	})
	return r
}

// func (b *Bot) SendDocument(arg *Body) *Result {}

// {
// 	"ok": false,
// 	"error_code": 400,
// 	"description": "[Error]: Bad Request: user not found"
// }

// SendMessage -
func (b *Bot) SendMessage(arg *Body) *Result {
	r, body := b.req("sendMessage", arg)

	if r.Status != OK {
		switch r.Desc {
		case "Bad Request: chat not found",
			"Forbidden: user is deactivated",
			"Forbidden: bot was blocked by the user",
			"Forbidden: bot can't send messages to bots":
			r.Status = BadChat
		case "Bad Request: message text is empty",
			"Bad Request: wrong parameter action in request",
			"Forbidden: bot was kicked from the group chat":
			r.Status = BadOther
		}

	} else {
		var res struct {
			ID int64 `json:"message_id"`
		}
		json.Unmarshal(body, &res)
		r.ID = res.ID
	}
	return r
}

func (b *Bot) getMe() *Result {
	r, _ := b.req("getMe", nil)
	return r
}

func (b *Bot) req(met string, arg any) (res *Result, by json.RawMessage) {
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
		return res, nil
	}
	res.Status = OK
	return res, body.Res
}
