package job

import (
	"fmt"
	"testing"

	"github.com/democracy-tools/countmein-density/internal"
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatapp"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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

func createDemonstration() error {

	dsc := ds.NewClientWrapper(env.Project)

	id, err := createDemonstrationInDatastore(dsc)
	if err != nil {
		return err
	}

	return inviteVolunteers(dsc, id)
}

func inviteVolunteers(dsc ds.Client, id string) error {

	var users []ds.User
	err := dsc.GetAll(ds.KindUser, &users)
	if err != nil {
		return err
	}

	wac := whatsapp.NewClientWrapper()
	log.Info("Sending invitations...")
	for _, currUser := range users {
		link := fmt.Sprintf("%s?user=%s&demonstration=%s", internal.JoinUrl, currUser.Id, id)
		log.Infof("%s (%s): %s\n%s", currUser.Name, currUser.Id, currUser.Phone, link)
		err = wac.SendInvitationTemplate(currUser.Phone, id, currUser.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func createDemonstrationInDatastore(dsc ds.Client) (string, error) {

	id := uuid.NewString()
	log.Infof("Creating demonstration '%s'...", id)
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

	var volunteers []ds.Volunteer
	err = dsc.GetFilter(ds.KindVolunteer, []ds.FilterField{{
		Name:     "demonstration_id",
		Operator: "=",
		Value:    demonstration.Id,
	}}, &volunteers)
	if err != nil {
		return err
	}

	logrus.Infof("Volunteer count: %d", len(volunteers))

	return nil
}
