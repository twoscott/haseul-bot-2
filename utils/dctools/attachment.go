package dctools

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

var attachmentClient = http.Client{
	Timeout: 30 * time.Second,
}

func DownloadAttachment(attachment discord.Attachment) ([]byte, error) {
	if attachment.URL == "" {
		return nil, errors.New("attachment must have a URL")
	}

	res, err := attachmentClient.Get(attachment.URL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err := errors.New(res.Status)
		log.Println("Attachment download error:", err)
		return nil, err
	}

	return io.ReadAll(res.Body)
}
