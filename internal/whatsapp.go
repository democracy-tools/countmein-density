package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/democracy-tools/countmein-density/internal/action"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/democracy-tools/go-common/slack"
	"github.com/sirupsen/logrus"
)

func (h *Handle) WhatsAppVerification(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query()
	mode := key.Get("hub.mode")
	token := key.Get("hub.verify_token")
	challenge := key.Get("hub.challenge")

	if mode == "subscribe" && token == h.whatsappVerificationToken {
		w.Write([]byte(challenge))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (h *Handle) WhatsAppEventHandler(w http.ResponseWriter, r *http.Request) {

	var payload whatsapp.WebhookMessage
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		logrus.Infof("failed to decode webhook message with '%v'", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	_ = forward(h.sc, payload)
	h.bot(payload)
}

func (h *Handle) bot(payload whatsapp.WebhookMessage) {

	if len(payload.Entry) == 1 && len(payload.Entry[0].Changes) == 1 {
		change := payload.Entry[0].Changes[0]
		if len(change.Value.Messages) == 1 && len(change.Value.Contacts) == 1 {
			request := action.Create(h.dsc, h.wac, change.Value.Contacts[0], change.Value.Messages[0])
			if request != nil {
				message, err := request.Run()
				if err != nil {
					h.sc.Debug(err.Error())
				} else if message != "" {
					h.sc.Info(message)
				}
			}
		}
	}
}

func forward(sc slack.Client, payload whatsapp.WebhookMessage) error {

	if len(payload.Entry) == 0 ||
		len(payload.Entry[0].Changes) == 0 ||
		len(payload.Entry[0].Changes[0].Value.Contacts) == 0 ||
		len(payload.Entry[0].Changes[0].Value.Messages) == 0 {
		logrus.Debug("ignoring whatsapp message with no contacts or no message")
		return nil
	}

	pretty, err := buildMessage(payload)
	if err != nil {
		return err
	}

	return sc.Info(string(pretty))
}

func buildMessage(message whatsapp.WebhookMessage) ([]byte, error) {

	if len(message.Entry) == 1 && len(message.Entry[0].Changes) == 1 {
		var res bytes.Buffer
		change := message.Entry[0].Changes[0]
		if len(change.Value.Contacts) == 1 {
			contact := change.Value.Contacts[0]
			res.WriteString(fmt.Sprintf("%s (%s)\n", contact.Profile.Name, contact.WaID))
			for _, currMessage := range change.Value.Messages {
				if currMessage.Type == whatsapp.TypeText {
					res.WriteString(fmt.Sprintf("%s\n", currMessage.Text.Body))
				} else if currMessage.Type == whatsapp.TypeButton {
					res.WriteString(fmt.Sprintf("%s\n", currMessage.Button.Text))
				} else {
					res.WriteString(fmt.Sprintf("%s\n", currMessage.Type))
				}
			}
			return res.Bytes(), nil
		}
	}

	pretty, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal whatsapp message with '%v'", err)
		return nil, err
	}
	return pretty, nil
}
