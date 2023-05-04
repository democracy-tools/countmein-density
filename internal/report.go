package internal

import (
	"fmt"
	"strings"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/slack"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/sirupsen/logrus"
)

// Send thanks to all volunteers
func report(dsc ds.Client, wac whatsapp.Client, sc slack.Client, from string, message string) error {

	err := validateUserAdmin(dsc, from, message)
	if err != nil {
		return err
	}

	template, count, url, err := getReportDetails(message)
	if err != nil {
		return err
	}

	demonstration, err := ds.GetKaplanDemonstration(dsc)
	if err != nil {
		return err
	}

	participantIdToPhone, err := getParticipants(dsc, sc, demonstration.Id)
	if err != nil {
		return err
	}

	for _, currPhone := range participantIdToPhone {
		err = wac.SendThanksTemplate(template, currPhone, url, []string{count})
		if err != nil {
			sc.Debug(err.Error())
		}
	}

	return nil
}

func validateUserAdmin(dsc ds.Client, from string, message string) error {

	ok, err := ds.IsAdmin(dsc, from)
	if err != nil {
		return err
	}
	if !ok {
		err = fmt.Errorf("user '%s' is not authorized to send report '%s'", from, message)
		logrus.Info(err.Error())
		return err
	}

	return nil
}

func getReportDetails(message string) (string, string, string, error) {

	split := strings.Split(message, " ")
	if len(split) == 3 && strings.HasPrefix(split[2], "https://") &&
		(strings.EqualFold(split[0], "thanks1") ||
			strings.EqualFold(split[0], "thanks4") ||
			strings.EqualFold(split[0], "thanks5") ||
			strings.EqualFold(split[0], "thanks6")) {
		return split[0], split[1], split[2], nil
	}

	err := fmt.Errorf("invalid report message '%s'", message)
	logrus.Error(err.Error())
	return "", "", "", err
}

func getParticipants(dsc ds.Client, sc slack.Client, demonstration string) (map[string]string, error) {

	// *** datastore does not support join and group-by ***

	var volunteers []ds.Volunteer
	err := dsc.GetFilter(ds.KindVolunteer, []ds.FilterField{{
		Name:     "demonstration_id",
		Operator: "=",
		Value:    demonstration,
	}}, &volunteers)
	if err != nil {
		return nil, err
	}

	var users []ds.User
	err = dsc.GetAll(ds.KindUser, &users)
	if err != nil {
		return nil, err
	}
	userIdToPhone := make(map[string]string)
	for _, currUser := range users {
		userIdToPhone[currUser.Id] = currUser.Phone
	}

	participantIdToPhone := make(map[string]string)
	for _, currVolunteer := range volunteers {
		currPhone, ok := userIdToPhone[currVolunteer.UserId]
		if ok {
			participantIdToPhone[currVolunteer.UserId] = currPhone
		} else {
			msg := fmt.Sprintf("user '%s' sent observation '%+v', but did not found in datastore user entity", currVolunteer.UserId, currVolunteer)
			logrus.Error(msg)
			sc.Debug(msg)
		}
	}

	return participantIdToPhone, nil
}
