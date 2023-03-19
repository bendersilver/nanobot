package nanobot

import (
	"context"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type MBot struct {
	ctxCancelFn context.CancelFunc
	cli         *telegram.Client
	sender      *message.Sender
}

func (m *MBot) Close() {
	m.ctxCancelFn()
}

func New(appID int, appHash, botToken, sessionFile string) (*MBot, error) {
	var m MBot

	m.cli = telegram.NewClient(appID, appHash, telegram.Options{
		SessionStorage: &session.FileStorage{
			Path: sessionFile,
		},
	})
	m.sender = message.NewSender(tg.NewClient(m.cli))
	e := make(chan error)
	var ctx context.Context
	ctx, m.ctxCancelFn = context.WithCancel(context.Background())
	go func(er chan error) {
		err := m.cli.Run(ctx, func(ctx context.Context) error {
			status, err := m.cli.Auth().Status(ctx)
			if err != nil {
				return err
			}
			if !status.Authorized {
				_, err = m.cli.Auth().Bot(ctx, botToken)
				if err != nil {
					return err
				}
			}

			err = m.cli.Ping(ctx)
			if err != nil {
				return err
			}

			close(er)
			<-ctx.Done()
			return nil
		})
		if err != nil {
			er <- err
		}
	}(e)

	return &m, <-e
}
