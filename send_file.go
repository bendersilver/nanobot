package nanobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type InputMedia string

const (
	InputDocument InputMedia = "document"
	InputAudio    InputMedia = "audio"
	InputPhoto    InputMedia = "photo"
	InputVideo    InputMedia = "video"
)

type FileUpload struct {
	File      io.Reader
	FileName  string
	Caption   string
	ParseMode ParseMode
}

type media struct {
	Type      InputMedia `json:"type"`
	Media     string     `json:"media"`
	Caption   string     `json:"caption,omitempty"`
	ParseMode ParseMode  `json:"parse_mode,omitempty"`
}

func uploadFiles(url string, chatID int64, fileType InputMedia, items ...*FileUpload) (b []byte, err error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	var multi bool = len(items) > 1
	if multi {
		url = url + "sendMediaGroup"
	} else {
		switch fileType {
		case InputDocument:
			url = url + "sendDocument"
		case InputAudio:
			url = url + "sendAudio"
		case InputPhoto:
			url = url + "sendPhoto"
		case InputVideo:
			url = url + "sendVideo"
		default:
			return nil, fmt.Errorf("not supported input media type")
		}
	}
	md := make([]*media, len(items))

	for ix, item := range items {
		var fieldname string
		fieldname = fmt.Sprintf("%s%d", fileType, ix)
		if !multi {
			fieldname = string(fileType)
		}
		fw, err := w.CreateFormFile(fieldname, item.FileName)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(fw, item.File)
		if err != nil {
			return nil, err
		}

		md[ix] = &media{
			Type:      fileType,
			Media:     "attach://" + fieldname,
			Caption:   item.Caption,
			ParseMode: item.ParseMode,
		}
		if !multi {
			err = w.WriteField("type", string(fileType))
			if err != nil {
				return nil, err
			}
			if item.Caption != "" {
				err = w.WriteField("caption", item.Caption)
				if err != nil {
					return nil, err
				}
				err = w.WriteField("parse_mode", item.ParseMode)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	if multi {
		bmedia, err := json.Marshal(md)
		if err != nil {
			return nil, err
		}
		err = w.WriteField("media", string(bmedia))
		if err != nil {
			return nil, err
		}
	}
	err = w.WriteField("chat_id", fmt.Sprint(chatID))
	if err != nil {
		return nil, err
	}
	w.Close()

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := new(http.Client).Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}
	return
}

func (b *Bot) SendDocuments(chatID int64, inputMedia InputMedia, fls ...*FileUpload) error {
	byt, err := uploadFiles(b.uri,
		chatID,
		inputMedia,
		fls...)
	if err != nil {
		return err
	}
	return fmt.Errorf("%s", byt)
}
