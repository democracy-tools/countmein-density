package action

import (
	"errors"
	"fmt"
	"strings"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/sirupsen/logrus"
)

type Report struct {
	dsc     ds.Client
	wac     whatsapp.Client
	from    string
	message string
}

func NewReport(dsc ds.Client, wac whatsapp.Client, phone string, message string) Request {

	return &Report{
		dsc:     dsc,
		wac:     wac,
		from:    phone,
		message: message,
	}
}

func (a *Report) Run() (string, error) {

	if err := report(a.dsc, a.wac, a.from, a.message); err != nil {
		return "", fmt.Errorf("%s failed to send report %s with %v", a.from, a.message, err)
	}
	return fmt.Sprintf("%s sent report %s", a.from, a.message), nil
}

func report(dsc ds.Client, wac whatsapp.Client, from string, message string) error {

	err := validateUserAdmin(dsc, from, message)
	if err != nil {
		return err
	}

	template, count, url, err := getReportDetails(message)
	if err != nil {
		return err
	}

	return sendReportToAllVolunteers(dsc, wac, template, from, url, count)
}

func sendReportToAllVolunteers(dsc ds.Client, wac whatsapp.Client,
	template, from, url, count string) error {

	demonstration, err := ds.GetKaplanDemonstration(dsc)
	if err != nil {
		return err
	}

	participantIdToPhone, err := getParticipants(dsc, demonstration.Id)
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
		err = fmt.Errorf("user '%s' is not authorized to send a report '%s'", from, message)
		logrus.Info(err.Error())
		return err
	}

	return nil
}

func getReportDetails(message string) (string, string, string, error) {

	split := strings.Split(message, " ")
	if isPublishCountOnly(split) {
		return strings.ToLower(split[0]), split[1], "", nil
	}
	if isPublishCountWithShare(split) {
		return strings.ToLower(split[0]), split[1], split[2], nil
	}

	err := fmt.Errorf("invalid report message '%s'", message)
	logrus.Error(err.Error())
	return "", "", "", err
}

func getParticipants(dsc ds.Client, demonstration string) (map[string]string, error) {

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
			logrus.Error(fmt.Sprintf("user '%s' sent observation '%+v', but did not found in datastore user entity", currVolunteer.UserId, currVolunteer))
		}
	}

	return participantIdToPhone, nil
}
