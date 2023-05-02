package infra

import (
	"testing"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestVolunteerCount(t *testing.T) {

	// env.Initialize()
	t.Skip("infra")
	require.NoError(t, countVolunteers())
}

func TestCreateDemonstration(t *testing.T) {

	// env.Initialize()
	t.Skip("infra")
	require.NoError(t, createDemonstration())
}

func TestInviteSpecificUser(t *testing.T) {

	// env.Initialize()
	t.Skip("infra")
	require.NoError(t, inviteSpecificUser("972123456789"))
}

func createDemonstration() error {

	dsc := ds.NewClientWrapper(env.Project)

	id, err := createDemonstrationInDatastore(dsc)
	if err != nil {
		return err
	}

	return inviteAllUsers(dsc, id)
}

func inviteSpecificUser(phone string) error {

	dsc := ds.NewClientWrapper(env.Project)

	demonstration, err := ds.GetKaplanDemonstration(dsc)
	if err != nil {
		return err
	}

	user, err := ds.GetUserByPhone(dsc, phone)
	if err != nil {
		return err
	}

	return inviteVolunteers(dsc, demonstration.Id, []ds.User{*user})
}

func inviteAllUsers(dsc ds.Client, demonstrationId string) error {

	var users []ds.User
	err := dsc.GetAll(ds.KindUser, &users)
	if err != nil {
		return err
	}

	return inviteVolunteers(dsc, demonstrationId, users)
}

func inviteVolunteers(dsc ds.Client, demonstrationId string, users []ds.User) error {

	wac := whatsapp.NewClientWrapper()
	logrus.Infof("Sending invitations... demonstration '%s'", demonstrationId)
	for _, currUser := range users {
		logrus.Infof("%s (%s)", currUser.Name, currUser.Phone)
		err := wac.SendInvitationTemplate(currUser.Phone)
		if err != nil {
			return err
		}
	}

	return nil
}

func createDemonstrationInDatastore(dsc ds.Client) (string, error) {

	id := uuid.NewString()
	logrus.Infof("Creating demonstration '%s'...", id)
	err := dsc.Put(ds.KindDemonstration, ds.DemonstrationKaplan,
		&ds.Demonstration{Id: id, Name: ds.DemonstrationKaplan})
	if err != nil {
		return "", err
	}

	return id, nil
}

func countVolunteers() error {

	dsc := ds.NewClientWrapper(env.Project)

	demonstration, err := ds.GetKaplanDemonstration(dsc)
	if err != nil {
		return err
	}

	volunteers, err := ds.GetVolunteers(dsc, demonstration.Id)
	if err != nil {
		return err
	}

	logrus.Infof("Volunteer count: %d", len(volunteers))

	return nil
}
