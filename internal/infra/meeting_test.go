package infra

import (
	"testing"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestSendMessageToAllUsers(t *testing.T) {

	// env.Initialize()
	t.Skip("infra")
	require.NoError(t, sendMessageToAllUsers())
}

func sendMessageToAllUsers() error {

	dsc := ds.NewClientWrapper(env.Project)

	var users []ds.User
	err := dsc.GetAll(ds.KindUser, &users)
	if err != nil {
		return err
	}

	wac := whatsapp.NewClientWrapper()
	for _, currUser := range users {
		logrus.Infof("%s (%s)", currUser.Name, currUser.Phone)
		err := wac.SendBodyParamsTemplate("meeting", currUser.Phone, []string{"חמישי", "20:30", "פידבק מההפגנות האחרונות ותכניות קדימה", "meet.google.com/fpf-bvbt-cvd"})
		if err != nil {
			return err
		}
	}

	return nil
}
