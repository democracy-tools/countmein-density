package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/democracy-tools/countmein-density/internal/env"
	log "github.com/sirupsen/logrus"
)

type MessageRequest struct {
	MessagingProduct string      `json:"messaging_product"`
	RecipientType    string      `json:"recipient_type"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Text             MessageText `json:"text"`
}

type MessageText struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type Client interface {
	Send(phone string, body string) error
	SendSignupTemplate(phone string, token string) error
}

type ClientWrapper struct {
	auth string
	from string
}

func NewClientWrapper() Client {
	return &ClientWrapper{
		auth: fmt.Sprintf("Bearer %s", env.GetWhatAppToken()),
		from: env.GetWhatsAppFromPhone(),
	}
}

func (c *ClientWrapper) SendSignupTemplate(to string, token string) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(TemplateMessage{
		To:   to,
		Type: "template",
		Template: Template{
			Namespace: "",
			Language: Language{
				Policy: "deterministic",
				Code:   "he",
			},
			Name: "signup2",
			Components: []Components{{
				Type:    "button",
				SubType: "url",
				Index:   "0",
				Parameters: []Parameters{{
					Type: "text",
					Text: token,
				}},
			}},
		},
	})
	if err != nil {
		log.Errorf("failed to encode whatsapp sigunup message request with '%v' target phone '%s'", err, to)
		return err
	}

	r, err := http.NewRequest(http.MethodPost, getMessageUrl(c.from), &buf)
	if err != nil {
		log.Errorf("failed to create HTTP request for sending a whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", c.auth)

	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		log.Errorf("failed to send whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		log.Infof("failed to send whatsapp message to '%s' with '%s'", to, response.Status)
		return err
	}

	return nil
}

func (c *ClientWrapper) Send(to string, message string) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(MessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "text",
		Text: MessageText{
			PreviewURL: false,
			Body:       message,
		},
	})
	if err != nil {
		log.Errorf("failed to encode whatsapp message request with '%v'. target phone '%s'", err, to)
		return err
	}

	r, err := http.NewRequest(http.MethodPost, getMessageUrl(c.from), &buf)
	if err != nil {
		log.Errorf("failed to create HTTP request for sending a whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", c.auth)

	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		log.Errorf("failed to send whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		log.Infof("failed to send whatsapp message to '%s' with '%s'", to, response.Status)
		return err
	}

	return nil
}

func getMessageUrl(from string) string {

	return fmt.Sprintf("https://graph.facebook.com/v16.0/%s/messages", from)
}
