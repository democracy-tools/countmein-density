package action_test

import (
	"testing"

	"github.com/democracy-tools/countmein-density/internal/action"
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAction_Run(t *testing.T) {

	message, err := action.Create(ds.NewInMemoryClient(),
		whatsapp.NewInMemoryClient(),
		whatsapp.Contact{
			WaID:    "123456789",
			Profile: whatsapp.ContactProfile{Name: "Israel"}},
		whatsapp.Message{
			Type: whatsapp.TypeText,
			Text: whatsapp.MessageText{Body: "קפלן"}}).Run()
	require.NoError(t, err)
	require.NotEmpty(t, message)
}
