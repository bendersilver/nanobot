# simple telegram bot api

```go
package main

import (
	"log"
	"strings"
	"time"

	"github.com/bendersilver/nanobot"
)

func main() {
	bot, err := nanobot.New("<token>")
	if err != nil {
		log.Fatal(err)
	}

	// SendMessage -
	rsp, err := bot.SendMessage("<chat_id>", "text <b>hehe</b>", nanobot.ModeHTML)
	if err != nil {
		log.Fatal(err)
	}
	switch rsp.Status {
	case nanobot.BadChat:
		// remove chat

	case nanobot.BadToken:
		// remove token

	case nanobot.Sleep:
		time.Sleep(time.Duration(rsp.Parms.RetryAfter) * time.Second)

	case nanobot.BadOther:
		log.Fatal(rsp.Desc)
	case nanobot.OK:
		// json.Unmarshal(rsp.Res, any)
	}

	// DeleteMessage
	rsp, err = bot.DeleteMessage("<chat_id>", "<message_id>")
	if err != nil {
		log.Fatal(err)
	}

	// Send file
	rsp, err = bot.SendDocuments("<chat_id>",
		nanobot.InputDocument,
		&nanobot.FileUpload{
			File:      strings.NewReader("one"),
			FileName:  "one.txt",
			Caption:   "<b>text</b>",
			ParseMode: "HTML",
		})
	if err != nil {
		log.Fatal(err)
	}

	// Send multiple file
	rsp, err = bot.SendDocuments("<chat_id>",
		nanobot.InputDocument,
		&nanobot.FileUpload{
			File:      strings.NewReader("one"),
			FileName:  "one.txt",
			Caption:   "<b>text</b>",
			ParseMode: "HTML",
		}, &nanobot.FileUpload{
			File:     strings.NewReader("two"),
			FileName: "two.txt",
			Caption:  "text",
		}, &nanobot.FileUpload{
			File:     strings.NewReader("three"),
			FileName: "three.txt",
		})
}

```