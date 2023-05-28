package action

import (
	"fmt"
	"strings"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/whatsapp"
	log "github.com/sirupsen/logrus"
)

var priority = []string{"A5", "A4", "A3", "A13", "A12", "A11", "A7", "A7B", "A8", "A10", "A9", "A18", "A19",
	"A20", "A24", "A23", "A14A", "A14", "A14B", "A17", "A17B", "K1", "K1A", "K2", "K2B", "K3", "K4", "K5", "K6", "K7", "K9", "K10", "K11",
	"K12", "K13", "K14", "K19", "K20", "K21", "K22", "K26"}

type Join struct {
	dsc     ds.Client
	wac     whatsapp.Client
	name    string
	phone   string
	message string
}

func NewJoin(dsc ds.Client, wac whatsapp.Client, phone string, name string) Request {

	return &Join{
		dsc:   dsc,
		wac:   wac,
		name:  name,
		phone: phone,
	}
}

func (a *Join) Run() (string, error) {

	res, err := join(a.dsc, a.wac, a.phone)
	if err != nil {
		return "", fmt.Errorf("user %s (%s) failed to join demonstration with %v", a.name, a.phone, err)
	}
	return res, nil
}

func join(dsc ds.Client, wac whatsapp.Client, phone string) (string, error) {

	user, err := ds.GetUserByPhone(dsc, phone)
	if err != nil {
		return "", err
	}

	demonstration, err := ds.GetKaplanDemonstration(dsc)
	if err != nil {
		return "", err
	}

	available, err := GetAvailablePolygons(dsc, demonstration.Id)
	if err != nil {
		return "", err
	}
	if len(available) == 0 {
		err = fmt.Errorf("no available polygon found, user '%s (%s)'", user.Name, user.Phone)
		log.Info(err.Error())
		return "", err
	}

	polygon, location := getPolygonByPriority(available, user.Preference)
	if polygon == "" {
		for polygon, location = range available {
			break
		}
	}

	err = dsc.Put(ds.KindVolunteer, ds.GetVolunteerId(demonstration.Id, user.Id), &ds.Volunteer{
		UserId:          user.Id,
		DemonstrationId: demonstration.Id,
		Polygon:         polygon,
		Location:        location,
		Time:            time.Now().Unix(),
	})
	if err != nil {
		return "", err
	}

	err = wac.SendDemonstrationTemplate(user.Phone, user.Id)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("volunteer added: %s (%s) polygon %s demonstration %s (%s)", user.Name, user.Phone, polygon, demonstration.Name, demonstration.Id), nil
}

func getPolygonByPriority(available map[string]string, preferred string) (string, string) {

	if preferred != "" {
		for _, curr := range strings.Split(strings.ReplaceAll(preferred, " ", ""), ",") {
			res, ok := available[curr]
			if ok {
				return curr, res
			}
		}
	}

	for _, currPolygon := range priority {
		res, ok := available[currPolygon]
		if ok {
			return currPolygon, res
		}
	}

	return "", ""
}
