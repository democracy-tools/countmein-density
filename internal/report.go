package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/democracy-tools/countmein-density/internal/ds"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/democracy-tools/go-common/slack"
	"github.com/sirupsen/logrus"
)

func ReportRequest(message []string) bool {

	return len(message) == 2 && strings.EqualFold(message[0], "thanks1")
}

func ReportRequestWithShare(message []string) bool {

	return len(message) == 3 &&
		strings.HasPrefix(message[2], "https://") &&
		(strings.EqualFold(message[0], "thanks4") ||
			strings.EqualFold(message[0], "thanks5") ||
			strings.EqualFold(message[0], "thanks6"))
}

func Report(dsc ds.Client, wac whatsapp.Client, sc slack.Client, from string, message string) error {

	err := validateUserAdmin(dsc, from, message)
	if err != nil {
		return err
	}

	return report(dsc, wac, sc, from, message)
}

func report(dsc ds.Client, wac whatsapp.Client, sc slack.Client, from string, message string) error {

	template, count, url, err := getReportDetails(message)
	if err != nil {
		return err
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

func getReportDetails(message string) (string, string, string, error) {

	split := strings.Split(message, " ")
	if ReportRequest(split) {
		return strings.ToLower(split[0]), split[1], "", nil
	}
	if ReportRequestWithShare(split) {
		return strings.ToLower(split[0]), split[1], split[2], nil
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
