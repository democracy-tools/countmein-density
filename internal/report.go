package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/slack"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/sirupsen/logrus"
)

const (
	AudienceAll = "all"
	AudienceMe  = "me"
)

func Report(dsc ds.Client, wac whatsapp.Client, sc slack.Client, from string, message string) error {

	err := validateUserAdmin(dsc, from, message)
	if err != nil {
		return err
	}

	return report(dsc, wac, sc, from, message)
}

func report(dsc ds.Client, wac whatsapp.Client, sc slack.Client, from string, message string) error {

	template, count, url, audience, err := getReportDetails(message)
	if err != nil {
		return err
	}

	if audience == AudienceMe {
		return sendReport(wac, []string{from}, template, url, count)
	}
	return sendReportToAllVolunteers(dsc, wac, sc, template, from, url, count)
}

func sendReportToAllVolunteers(dsc ds.Client, wac whatsapp.Client, sc slack.Client,
	template, from, url, count string) error {

	demonstration, err := ds.GetKaplanDemonstration(dsc)
	if err != nil {
		return err
	}

	participantIdToPhone, err := getParticipants(dsc, sc, demonstration.Id)
	if err != nil {
		return err
	}

	return sendReport(wac, getPhones(participantIdToPhone), template, url, count)
}

func sendReport(wac whatsapp.Client, phones []string, template string, url string, count string) error {

	var errs error
	for _, currPhone := range phones {
		err := wac.SendThanksTemplate(template, currPhone, url, []string{count})
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if errs != nil {
		return errs
	}

	return nil
}

func getPhones(participantIdToPhone map[string]string) []string {

	var res []string
	for _, currPhone := range participantIdToPhone {
		res = append(res, currPhone)
	}

	return res
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

func getReportDetails(message string) (string, string, string, string, error) {

	split := strings.Split(message, " ")
	if (len(split) == 3 || len(split) == 4) && strings.HasPrefix(split[2], "https://") {
		template := strings.ToLower(split[0])
		if template == "thanks1" || template == "thanks4" ||
			template == "thanks5" || template == "thanks6" {
			return template, split[1], split[2], getAudience(split), nil
		}
	}

	err := fmt.Errorf("invalid report message '%s'", message)
	logrus.Error(err.Error())
	return "", "", "", "", err
}

func getAudience(message []string) string {

	if len(message) == 4 && message[3] == AudienceMe {
		return AudienceMe
	}
	return AudienceAll
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
