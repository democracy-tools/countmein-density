package infra

import (
	"testing"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/stretchr/testify/require"
)

func TestSendWhatsAppMessageToNewUsers(t *testing.T) {

	t.Skip("infra")
	// env.Initialize()

	// 6 days ago
	// fmt.Println(time.Now().Add(time.Hour * 24 * (-6)).Unix())

	dsc := ds.NewClientWrapper(env.Project)
	var users []ds.User
	require.NoError(t, dsc.GetFilter(ds.KindUser,
		[]ds.FilterField{{
			Name:     "time",
			Operator: ">",
			Value:    1682853520,
		}}, &users))

	wac := whatsapp.NewClientWrapper()
	for _, curr := range users {
		require.NoError(t, wac.Send(curr.Phone, "היי תודה שהצטרפת אלינו! רצינו לעדכן שניתן לשנות את הפוליגון שלך (האזור שהוקצה לך לדווח ממנו) ע״י הלינק בתחתית העמוד של האפליקציה - שינוי פוליגון. נתראה בהפגנה :)"))
	}
}
