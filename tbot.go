package nanobot

import (
	"encoding/json"
	"fmt"
	"io"
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
	// Sleep -
	Sleep
	// BadOther -
	BadOther
)
const tURL = "https://api.telegram.org/bot%s/"

type Result struct {
	Status Bad
	OK     bool            `json:"ok"`
	Code   int             `json:"error_code"`
	Desc   string          `json:"description"`
	Res    json.RawMessage `json:"result"`
	Parms  struct {
		MtoChatID  int64 `json:"migrate_to_chat_id"`
		RetryAfter int   `json:"retry_after"`
	} `json:"parameters"`
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
	r, err := b.GetMe()
	if err != nil {
		return nil, err
	}
	if !r.OK {
		return nil, fmt.Errorf(r.Desc)
	}
	return &b, nil
}

func chekOK(rsp *http.Response) (res *Result) {
	defer rsp.Body.Close()
	res = new(Result)
	res.Code = rsp.StatusCode
	res.Status = BadOther

	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		res.Desc = "cannot read response body"
		return
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		res.Desc = fmt.Sprintf("%v\n%s", err, b)
		return
	}
	if !res.OK {
		switch res.Code {
		case 421:
			res.Status = BadToken
		case 429:
			res.Status = Sleep
		default:
			switch res.Desc {
			case "Bad Request: chat not found",
				"Forbidden: user is deactivated",
				"Forbidden: bot was blocked by the user",
				"Forbidden: bot can't send messages to bots",
				"Forbidden: bot was kicked from the group chat",
				"Bad Request: user not found":
				res.Status = BadChat
			}
		}
	}
	return
}
