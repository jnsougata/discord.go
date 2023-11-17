package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

var DISCORD_API_VERSION = "10"
var DISCORD_API_BASE_URL = "https://discord.com/api/v" + DISCORD_API_VERSION

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func FormData(data any, files []File) ([]byte, string) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	payload, _ := json.Marshal(data)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="payload_json"`)
	header.Set("Content-Type", `application/json`)
	field, _ := writer.CreatePart(header)
	_, _ = field.Write(payload)
	for i, f := range files {
		ff, _ := writer.CreateFormFile(fmt.Sprintf(`files[%v]`, i), quoteEscaper.Replace(f.Name))
		_, _ = ff.Write(f.Content)
	}
	_ = writer.Close()
	return buffer.Bytes(), writer.Boundary()
}


type RequestOptions struct {
	Method    string
	Path      string
	Authorize bool
	Body      io.Reader
	Boundary  string
	Kwargs    map[string]interface{}
}

type HttpClient struct {
	state *AppState
}

func (c *HttpClient) Request(o RequestOptions) (*http.Response, error) {
	url := DISCORD_API_BASE_URL + o.Path
	req, err := http.NewRequest(o.Method, url, o.Body)
	if err != nil {
		return nil, err
	}
	if o.Boundary != "" {
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+o.Boundary)
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	if o.Authorize {
		req.Header.Set("Authorization", "Bot "+c.state.DiscordToken)
	}
	if o.Kwargs != nil {
		if reason, ok := o.Kwargs["reason"]; ok {
			req.Header.Set("X-Audit-Log-Reason", reason.(string))
		}
	}
	return http.DefaultClient.Do(req)
}

func (c *HttpClient) sync(commands []ApplicationCommand) (*http.Response, error) {
	return c.Request(
		RequestOptions{
			Method:    http.MethodPut,
			Path:      fmt.Sprintf("/applications/%s/commands", c.state.ApplicationId),
			Authorize: true,
			Body:      ReaderFromAny(commands),
		})
}

func (c *HttpClient) SendInteractionCallback(
	interaction *Interaction,
	kind InteractionCallbackType,
	payload SendingPayload,
) (*http.Response, error) {
	f := map[string]interface{}{"type": int(kind), "data": payload}
	data, boundary := FormData(f, payload.Attchments)
	return c.Request(RequestOptions{
		Method:    http.MethodPost,
		Path:      fmt.Sprintf("/interactions/%s/%s/callback", interaction.Id, interaction.Token),
		Authorize: true,
		Body:      bytes.NewReader(data),
		Boundary:  boundary,
	})
}

func NewHttpClient(state *AppState) *HttpClient {
	return &HttpClient{state}
}
