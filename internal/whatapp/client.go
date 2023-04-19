package whatsapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/democracy-tools/countmein-density/internal/env"
	log "github.com/sirupsen/logrus"
)

type Client interface {
	Send(phone string, body string) error
	SendSignupTemplate(phone string, token string) error
	SendVerifyTemplate(phone string) error
	SendInvitationTemplate(to string, demonstration string, userId string) error
	SendDemonstrationTemplate(to string, demonstration string, userId string,
		user string, polygon string, location string) error
	SendThanksTemplate(to string, count string) error
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
	err := json.NewEncoder(&buf).Encode(newTemplate("signup", to, token, nil))
	if err != nil {
		log.Errorf("failed to encode whatsapp sigunup message request with '%v' phone '%s'", err, to)
		return err
	}

	return send(c.from, to, &buf, c.auth)
}

func (c *ClientWrapper) SendVerifyTemplate(to string) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(newTemplate("verify5", to, "", nil))
	if err != nil {
		log.Errorf("failed to encode whatsapp verify message request with '%v' phone '%s'", err, to)
		return err
	}

	return send(c.from, to, &buf, c.auth)
}

func (c *ClientWrapper) SendInvitationTemplate(to string, demonstration string, userId string) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(newTemplate("attend", to,
		fmt.Sprintf("?demonstration=%s&user=%s", demonstration, userId), nil))
	if err != nil {
		log.Errorf("failed to encode whatsapp attend message request with '%v' phone '%s'", err, to)
		return err
	}

	return send(c.from, to, &buf, c.auth)
}

func (c *ClientWrapper) SendDemonstrationTemplate(to string, demonstration string, userId string,
	user string, polygon string, location string) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(newTemplate("demonstration", to,
		fmt.Sprintf("?demonstration=%s&user-id=%s&user=%s&polygon=%s&q=%s",
			demonstration, userId, user, polygon, location), nil))
	if err != nil {
		log.Errorf("failed to encode whatsapp demonstration message request with '%v' user '%s (%s)'", err, user, userId)
		return err
	}

	return send(c.from, to, &buf, c.auth)
}

func (c *ClientWrapper) SendThanksTemplate(to string, count string) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(newTemplate("thanks", to, "", []string{count}))
	if err != nil {
		log.Errorf("failed to encode whatsapp thanks message request with '%v' phone '%s'", err, to)
		return err
	}

	return send(c.from, to, &buf, c.auth)
}

func (c *ClientWrapper) Send(to string, message string) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(TextMessageRequest{
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

	return send(c.from, to, &buf, c.auth)
}

func send(from string, to string, body io.Reader, auth string) error {

	r, err := http.NewRequest(http.MethodPost, getMessageUrl(from), body)
	if err != nil {
		log.Errorf("failed to create HTTP request for sending a whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", auth)

	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		log.Errorf("failed to send whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		msg := fmt.Sprintf("failed to send whatsapp message to '%s' with '%s'", to, response.Status)
		log.Info(msg)
		return errors.New(msg)
	}

	return nil
}

func newTemplate(name string, to string, buttonUrlParam string, bodyTextParams []string) TemplateMessageRequest {

	res := TemplateMessageRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: Template{
			Name: name,
			Language: Language{
				Policy: "deterministic",
				Code:   "he",
			},
		},
	}

	if buttonUrlParam != "" {
		res.Template.Components = []Component{{
			Type:    "button",
			SubType: "url",
			Index:   "0",
			Parameters: []Parameter{{
				Type: "text",
				Text: buttonUrlParam,
			}},
		}}
	}

	for _, currParam := range bodyTextParams {
		res.Template.Components = append(res.Template.Components, Component{
			Type: "body",
			Parameters: []Parameter{{
				Type: "text",
				Text: currParam,
			}}})
	}

	return res
}

func getMessageUrl(from string) string {

	return fmt.Sprintf("https://graph.facebook.com/v16.0/%s/messages", from)
}
