package action_test

import (
	"testing"

	"github.com/democracy-tools/countmein-density/internal/action"
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/stretchr/testify/require"
)

func TestSendRegretInvitation_Run(t *testing.T) {

	message, err := action.Create(ds.NewInMemoryClient(),
		whatsapp.NewInMemoryClient(),
		whatsapp.Contact{
			WaID:    "123445566",
			Profile: whatsapp.ContactProfile{Name: "Israel"}},
		whatsapp.Message{
			Type:   whatsapp.TypeButton,
			Button: whatsapp.Button{Text: "לא הפעם"}}).Run()
	require.NoError(t, err)
	require.Empty(t, message)
}
