package internal

import (
	"encoding/json"
	"net/http"

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
		logrus.Infof("failed to decode webhook message with '%v'")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(payload.Entry) == 1 && len(payload.Entry[0].Changes) == 1 {
		change := payload.Entry[0].Changes[0]
		if len(change.Value.Messages) == 1 {
			message := change.Value.Messages[0]
			if message.Type == "text" && message.Text.Body == "join" {
				contact := change.Value.Contacts[0]
				createUser(h.dsc, h.wac, contact.WaID, contact.Profile.Name, "")
			}
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
