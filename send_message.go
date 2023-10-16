package nanobot

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type message struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int64  `json:"message_id,omitempty"`
	Text      string `json:"text,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func (b *Bot) request(path string, id, mID int64, msg string, mode ParseMode) (*Result, error) {
	data, err := json.Marshal(&message{
		ChatID:    id,
		MessageID: mID,
		Text:      msg,
		ParseMode: mode,
	})
	if err != nil {
		return nil, err
	}
	rsp, err := http.Post(b.uri+path, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return chekOK(rsp), nil
}

// SendMessage -
func (b *Bot) SendMessage(chatID int64, message string, parseMode ParseMode) (*Result, error) {
	return b.request("sendMessage", chatID, 0, message, parseMode)
}

// DeleteMessage -
func (b *Bot) DeleteMessage(chatID, messageID int64) (*Result, error) {
	return b.request("deleteMessage", chatID, messageID, "", "")
}

// GetMe -
func (b *Bot) GetMe() (*Result, error) {
	return b.request("getMe", 0, 0, "", "")
}
