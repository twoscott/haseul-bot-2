package htmlutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/utils/httputil"
)

type TemplateBody struct {
	TemplateHtml []byte          `json:"templateHtml,omitempty"`
	TemplateName string          `json:"templateName,omitempty"`
	Width        int             `json:"width,omitempty"`
	Height       int             `json:"height,omitempty"`
	Format       string          `json:"imageFormat,omitempty"`
	Quality      int             `json:"quality,omitempty"`
	Context      json.RawMessage `json:"context"`
}

func getSushiiImageServerURL() string {
	sushiiCfg := config.GetInstance().SushiiImageServer
	return fmt.Sprintf("http://%s:%d", sushiiCfg.Host, sushiiCfg.Port)
}

func getTemplateURL() string {
	return getSushiiImageServerURL() + "/template"
}

func templateToImage(
	templateName string,
	jsonContext []byte,
	width, height int,
	format string,
	quality int) ([]byte, error) {

	url := getTemplateURL()

	body := TemplateBody{
		TemplateName: templateName,
		Width:        width,
		Height:       height,
		Format:       format,
		Quality:      quality,
		Context:      jsonContext,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	jsonBuf := bytes.NewBuffer(bodyJson)
	res, err := httputil.Post(
		url,
		httputil.ContentTypeApplicationJSON,
		jsonBuf,
	)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf(
			"error fetching data from Sushii Image Server: %s",
			errors.New(res.Status),
		)
		return nil, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// TemplateToJPEG returns a JPEG of the rendered html corresponding to the
// provided HTML template, with the provided context data inserted.
func TemplateToJPEG(
	templateName string,
	jsonContext []byte,
	width, height int,
	quality int) ([]byte, error) {

	return templateToImage(
		templateName, jsonContext, width, height, "jpeg", quality,
	)
}

// TemplateToJPEG returns a PNG of the rendered html corresponding to the
// provided HTML template, with the provided context data inserted.
func TemplateToPNG(
	templateName string,
	jsonContext []byte,
	width, height int) ([]byte, error) {

	return templateToImage(templateName, jsonContext, width, height, "png", 0)
}
