package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/democracy-tools/countmein-density/internal/slack"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatapp"
	"github.com/sirupsen/logrus"
)

func (h *Handle) WhatsAppVerification(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query()
	mode := key.Get("hub.mode")
	token := key.Get("hub.verify_token")
	challenge := key.Get("hub.challenge")

	if len(mode) > 0 && len(token) > 0 {
		if mode == "subscribe" && token == h.whatsappVerificationToken {
			w.Write([]byte(challenge))
			return
		}
		w.WriteHeader(http.StatusForbidden)
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

	forward(h.sc, payload)

	if len(payload.Entry) == 1 && len(payload.Entry[0].Changes) == 1 {
		change := payload.Entry[0].Changes[0]
		if len(change.Value.Messages) == 1 && len(change.Value.Contacts) == 1 {
			contact := change.Value.Contacts[0]
			message := change.Value.Messages[0]
			if message.Type == whatsapp.TypeText {
				if isRegisterRequest(message.Text.Body) {
					code := createUser(h.dsc, h.wac, contact.WaID, contact.Profile.Name, "")
					if code == http.StatusCreated {
						h.sc.Info(fmt.Sprintf("User added: %s (%s)", contact.Profile.Name, contact.WaID))
					} else {
						h.sc.Debug(fmt.Sprintf("Failed to add user %s (%s) with %d", contact.Profile.Name, contact.WaID, code))
					}
				} else if isReportRequest(message.Text.Body) {
					err := report(h.dsc, h.wac, h.sc, contact.WaID, message.Text.Body)
					if err != nil {
						h.sc.Debug(fmt.Sprintf("Failed to send report %s with %v", message.Text.Body, err))
					} else {
						h.sc.Debug(fmt.Sprintf("Sent report %s", message.Text.Body))
					}
				}
			} else if message.Type == whatsapp.TypeButton {
				if isUnsubscribeRequestButton(message.Button.Text) {
					err := deleteUser(h.dsc, contact.WaID)
					if err != nil {
						h.sc.Info(fmt.Sprintf("Failed to delete user %s (%s) with %v", contact.Profile.Name, contact.WaID, err))
					} else {
						h.sc.Info(fmt.Sprintf("User deleted: %s (%s)", contact.Profile.Name, contact.WaID))
					}
				} else if isJoinRequestButton(message.Button.Text) {
					err := h.join(contact.WaID)
					if err != nil {
						h.sc.Info(fmt.Sprintf("User %s (%s) failed to join demonstration with %v", contact.Profile.Name, contact.WaID, err))
					}
				}
			}
		}
	}

	w.WriteHeader(http.StatusAccepted)
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

func isUnsubscribeRequestButton(message string) bool {

	return message == "בבקשה להסיר אותי"
}

func isRegisterRequest(message string) bool {

	message = strings.ReplaceAll(message, " ", "")
	return strings.EqualFold(message, "join") || message == "קפלן" ||
		message == "אנירוצהלהתנדבבספירתהמפגיניםבקפלן"
}

func isReportRequest(message string) bool {

	split := strings.Split(message, " ")
	return len(split) == 2 &&
		(strings.EqualFold(split[0], "thanks1") || strings.EqualFold(split[0], "thanks2"))
}

func isJoinRequestButton(message string) bool {

	return message == "כן, אני בעניין"
}
