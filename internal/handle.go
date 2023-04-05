package internal

import (
	"github.com/democracy-tools/countmein-density/internal/ds"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatapp"
)

type Handle struct {
	dsc ds.Client
	wac whatsapp.Client
}

func NewHandle(dsc ds.Client, wac whatsapp.Client) *Handle {

	return &Handle{dsc: dsc, wac: wac}
}
