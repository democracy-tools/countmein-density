package internal

import (
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	"github.com/democracy-tools/countmein-density/internal/slack"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatapp"
)

type Handle struct {
	dsc                       ds.Client
	wac                       whatsapp.Client
	sc                        slack.Client
	whatsappVerificationToken string
}

func NewHandle(dsc ds.Client, wac whatsapp.Client) *Handle {

	return &Handle{
		dsc:                       dsc,
		wac:                       wac,
		sc:                        slack.NewClientWrapper(),
		whatsappVerificationToken: env.GetWhatsAppVerificationToken(),
	}
}
