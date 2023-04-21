package job

import (
	"fmt"
	"strings"
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

func TestChangePolygon(t *testing.T) {

	// env.Initialize()
	t.Skip("infra")
	require.NoError(t, changePolygon("972501234567", "A78"))
}

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

func changePolygon(phone string, polygon string) error {

	dsc := ds.NewClientWrapper(env.Project)

	user, err := ds.GetUserByPhone(dsc, phone)
	if err != nil {
		return err
	}

	if sliceContains(strings.Split(strings.ReplaceAll(user.Preference, " ", ""), ","), polygon) {
		return fmt.Errorf("%s (%s) asked to change into polygon '%s', but has it as part of preference '%s'", user.Name, user.Phone, polygon, user.Preference)
	}
	err = dsc.Put(ds.KindUser, user.Id, &ds.User{
		Id:         user.Id,
		Phone:      user.Phone,
		Name:       user.Name,
		Preference: getPreference(user.Preference, polygon),
		Time:       user.Time,
		Role:       user.Role,
	})
	if err != nil {
		return err
	}

	return internal.NewHandle(dsc, whatsapp.NewClientWrapper()).Join(user)
}

func getPreference(preference string, polygon string) string {

	if preference == "" {
		return polygon
	}

	return fmt.Sprintf("%s,%s", preference, polygon)
}

func sliceContains(slice []string, item string) bool {

	for _, curr := range slice {
		if curr == item {
			return true
		}
	}

	return false
}
