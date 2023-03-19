package nanobot

import (
	"context"

	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/tg"
)

func (m *MBot) DeleteMessages(ids ...int) (bool, error) {
	a, err := tg.NewClient(m.cli).MessagesDeleteMessages(context.Background(), &tg.MessagesDeleteMessagesRequest{
		Revoke: true,
		ID:     ids,
	})
	if err != nil {
		return false, err
	}
	return a.PtsCount > 0, nil
}

func (m *MBot) SendText(userID int64, msg string) (int, error) {
	u, err := m.sender.To(&tg.InputPeerUser{
		UserID: userID,
	}).Text(context.Background(), msg)
	if err != nil {
		return 0, err
	}
	if sm, ok := u.(*tg.UpdateShortSentMessage); ok {
		return sm.ID, nil
	}
	return 0, nil
}

func (m *MBot) SendHTML(userID int64, msg string) (int, error) {
	u, err := m.sender.To(&tg.InputPeerUser{
		UserID: userID,
	}).StyledText(context.Background(), html.String(nil, msg))
	if err != nil {
		return 0, err
	}
	if sm, ok := u.(*tg.UpdateShortSentMessage); ok {
		return sm.ID, nil
	}
	return 0, nil
}
