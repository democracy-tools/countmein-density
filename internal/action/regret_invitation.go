package action

import (
	"fmt"

	"github.com/democracy-tools/countmein-density/internal/whatsapp"
)

type SendRegretInvitation struct {
	wac   whatsapp.Client
	phone string
}

func NewSendRegretInvitation(wac whatsapp.Client, phone string) Request {

	return &SendRegretInvitation{wac: wac, phone: phone}
}

func (a *SendRegretInvitation) Run() (string, error) {

	err := a.wac.SendRegretInvitationTemplate(a.phone)
	if err != nil {
		return "", fmt.Errorf("failed to send regret invitation to %s with %v", a.phone, err)
	}

	return "", nil
}
