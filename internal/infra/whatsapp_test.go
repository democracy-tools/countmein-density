package infra

import (
	"testing"

	whatsapp "github.com/democracy-tools/countmein-density/internal/whatapp"
	"github.com/stretchr/testify/require"
)

func TestSendWhatsAppMessage(t *testing.T) {

	// env.Initialize()
	t.Skip("infra")

	wac := whatsapp.NewClientWrapper()
	require.NoError(t, wac.Send("9721234567", "בהצלחה :)"))
}
