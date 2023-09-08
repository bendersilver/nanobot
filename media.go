package nanobot

import (
	"context"
	"io"

	"github.com/bendersilver/jlog"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

func (m *MBot) Upload(name string, rd io.ReadSeeker) (tg.InputFileClass, error) {

	mtype, err := mimetype.DetectReader(rd)
	if err != nil {
		return nil, err
	}
	_, err = rd.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	jlog.Notice(mtype)

	up := uploader.NewUploader(tg.NewClient(m.cli))
	fl, err := up.FromReader(context.Background(), name, rd)
	if err != nil {
		return nil, err
	}

	jlog.Noticef("%T", fl)
	return fl, nil
}

// func (m *MBot) SendPhoto(userID int64, f tg.InputFileClass, msg string) error {
// 	// m.sender.WithUploader(message.UploadedPhoto(f))
// 	u, err := m.sender.To(&tg.InputPeerUser{
// 		UserID: userID,
// 	}).UploadMedia(context.Background(), message.UploadedPhoto(f, html.String(nil, msg))) //.Media(context.Background(), message.Photo(f, ))
// 	// .Photo(context.Background(), &tg.InputMediaUploadedPhoto{}, html.String(nil, msg))
// 	if err != nil {
// 		return err
// 	}
// 	r, err := m.sender.To(&tg.InputPeerUser{
// 		UserID: userID,
// 	}).Media(context.Background(), message.Photo())
// 	jlog.Notice(u)
// 	return nil
// }

func (m *MBot) SendVideo(userID int64, f tg.InputFileClass, msg string) error {

	u, err := m.sender.To(&tg.InputPeerUser{
		UserID: userID,
	}).Video(context.Background(), f, html.String(nil, msg))
	if err != nil {
		return err
	}
	jlog.Notice(u)
	return nil
}
