package nanobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type FileUpload struct {
	File      io.Reader
	FileName  string
	Caption   string
	ParseMode string
}

func UploadItem(url, fileType string, chatID int64, item *FileUpload) (b []byte, err error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if x, ok := item.File.(io.Closer); ok {
		defer x.Close()
	}
	fw, err := w.CreateFormFile(fileType, item.FileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, item.File)
	if err != nil {
		return nil, err
	}

	err = w.WriteField("type", fileType)
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

func Upload(url, tpe string, values ...FileUpload) (err error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	// var media []map[string]string
	// if len(values) > 1 {
	media := make([]map[string]string, len(values))
	// }

	for ix, item := range values {
		if x, ok := item.File.(io.Closer); ok {
			defer x.Close()
		}
		fieldname := func() string {
			if len(values) > 1 {
				return fmt.Sprintf("%s%d", tpe, ix)
			}
			return tpe
		}()
		fw, err := w.CreateFormFile(fieldname, item.FileName)
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, item.File)
		if err != nil {
			return err
		}
		if len(media) > 1 {
			media[ix] = map[string]string{
				"type":       tpe,
				"media":      "attach://" + fieldname,
				"caption":    item.Caption,
				"parse_mode": "HTML",
			}
		} else {
			err = w.WriteField("type", tpe)
			if err != nil {
				return err
			}
			err = w.WriteField("caption", item.Caption)
			if err != nil {
				return err
			}
		}

	}
	bmedia, err := json.Marshal(&media)
	if err != nil {
		return err
	}
	err = w.WriteField("media", string(bmedia))
	if err != nil {
		return err
	}
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := new(http.Client).Do(req)
	if err != nil {
		return
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	res.Body.Close()
	// res.Body
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s\n%s", res.Status, b)
	}
	return
}

func (b *Bot) SendDocument(chatID int64, f *FileUpload) error {
	byt, err := UploadItem(
		b.uri+"sendDocument",
		"document",
		chatID, f)
	if err != nil {
		return err
	}
	return fmt.Errorf("%s", byt)
}

func (b *Bot) SendDocumentsd(arg *Body) error {
	return Upload(b.uri+"sendMediaGroup?chat_id=-1001695800269&parse_mode=HTML", "document", FileUpload{
		File:     strings.NewReader("text here"),
		FileName: "doc.txt",
		Caption:  "<b>text</b> <i>test</i>",
	},
		FileUpload{
			File:     strings.NewReader("text 2 here"),
			FileName: "doc2.txt",
			Caption:  "<b>text 2</b> <i>test 2</i>",
		})
}
