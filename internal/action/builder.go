package action

import (
	"strings"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/whatsapp"
)

func Create(dsc ds.Client, wac whatsapp.Client,
	contact whatsapp.Contact, message whatsapp.Message) Request {

	if message.Type == whatsapp.TypeText {
		if isRegister(message.Text.Body) {
			return NewCreateUser(dsc, wac, contact.WaID, contact.Profile.Name)
		} else if isPublishCount(message.Text.Body) {
			return NewReport(dsc, wac, contact.WaID, message.Text.Body)
		}
	} else if message.Type == whatsapp.TypeButton {
		if isJoinRequestButton(message.Button.Text) {
			return NewJoin(dsc, wac, contact.WaID, contact.Profile.Name)
		} else if isJoinNotThisTimeButton(message.Button.Text) {
			return NewSendRegretInvitation(wac, contact.WaID)
		}
	}

	return nil
}

func isRegister(message string) bool {

	return message == "קפלן" ||
		message == "אני רוצה להתנדב בספירת המפגינים בקפלן"
}

func isPublishCount(message string) bool {

	split := strings.Split(message, " ")
	return isPublishCountOnly(split) || isPublishCountWithShare(split)
}

func isPublishCountOnly(message []string) bool {

	return len(message) == 2 &&
		(strings.EqualFold(message[0], "thanks1") ||
			strings.EqualFold(message[0], "thanks11") ||
			strings.EqualFold(message[0], "thanks3"))
}

func isPublishCountWithShare(message []string) bool {

	return len(message) == 3 &&
		strings.HasPrefix(message[2], "https://") &&
		(strings.EqualFold(message[0], "thanks4") ||
			strings.EqualFold(message[0], "thanks5") ||
			strings.EqualFold(message[0], "thanks6"))
}

func isJoinRequestButton(message string) bool {

	return message == "כן, אני בעניין" || message == "להצטרפות"
}

func isJoinNotThisTimeButton(message string) bool {

	return message == "לא הפעם"
}
