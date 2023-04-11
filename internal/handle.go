package internal

import (
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatapp"
)

type Handle struct {
	dsc                       ds.Client
	wac                       whatsapp.Client
	whatsappVerificationToken string
}

func NewHandle(dsc ds.Client, wac whatsapp.Client) *Handle {

	return &Handle{dsc: dsc, wac: wac, whatsappVerificationToken: env.GetWhatsAppVerificationToken()}
}
