package infra

import (
	"fmt"
	"testing"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	"github.com/stretchr/testify/require"
)

func TestReport_VolunteerVsObservation(t *testing.T) {

	// env.Initialize()
	t.Skip("infra")
	require.NoError(t, volunteerVsObservation())
}

func volunteerVsObservation() error {

	dsc := ds.NewClientWrapper(env.Project)

	demonstration, err := ds.GetKaplanDemonstration(dsc)
	if err != nil {
		return err
	}

	observations, err := ds.GetObservations(dsc, demonstration.Id)
	if err != nil {
		return err
	}
	userIdToObservations := createUserIdToObservations(observations)

	volunteers, err := ds.GetVolunteers(dsc, demonstration.Id)
	if err != nil {
		return err
	}

	var userIds []string
	for _, currVolunteer := range volunteers {
		if _, ok := userIdToObservations[currVolunteer.UserId]; !ok {
			userIds = append(userIds, currVolunteer.UserId)
		}
	}

	for _, currUserId := range userIds {
		var user ds.User
		err = dsc.Get(ds.KindUser, currUserId, &user)
		if err != nil {
			return err
		}
		fmt.Printf("%s (%s)\n", user.Name, user.Phone)
	}

	return nil
}

func createUserIdToObservations(observations []ds.Observation) map[string][]ds.Observation {

	res := make(map[string][]ds.Observation)
	for _, currObservation := range observations {
		res[currObservation.User] = append(res[currObservation.User], currObservation)
	}

	return res
}
