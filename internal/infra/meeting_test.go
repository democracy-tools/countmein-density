package job

import (
	"testing"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
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

	// wac := whatsapp.NewClientWrapper()
	// for _, currUser := range users {
	//  log.Infof()
	// 	err = wac.SendTemplate(currUser.Phone, id, currUser.Id)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
