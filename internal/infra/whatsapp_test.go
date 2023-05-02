package infra

import (
	"testing"

	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/stretchr/testify/require"
)

func TestSendWhatsAppMessage(t *testing.T) {

	// env.Initialize()
	t.Skip("infra")

	wac := whatsapp.NewClientWrapper()
	require.NoError(t, wac.Send("9721234567", "היי, רצינו לעדכן שניתן לשנות מיקום ע״י הלינק באפליקציה של שינוי פוליגון. נתראה בהפגנה :)"))
}
