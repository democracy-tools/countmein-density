package internal

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	log "github.com/sirupsen/logrus"
)

var priority = []string{"A5", "A4", "A3", "A13", "A12", "A11", "A7", "A8", "A10", "A9", "A18", "A19",
	"A20", "A24", "A23", "A14", "A17", "K1", "K2", "K3", "K4", "K5", "K6", "K7", "K9", "K10", "K11",
	"K12", "K13", "K14", "K19", "K20", "K21", "K22", "K26"}

func (h *Handle) Join(user *ds.User) error {

	demonstration, err := ds.GetKaplanDemonstration(h.dsc)
	if err != nil {
		return err
	}

	available, err := getAvailablePolygons(h.dsc, demonstration.Id)
	if err != nil {
		return err
	}
	if len(available) == 0 {
		err = fmt.Errorf("no available polygon found for '%s: %s'", user.Name, user.Phone)
		log.Info(err.Error())
		return err
	}

	polygon, location := getPolygonByPriority(available, user.Preference)
	if polygon == "" {
		for polygon, location = range available {
			break
		}
	}

	err = h.dsc.Put(ds.KindVolunteer, ds.GetVolunteerId(demonstration.Id, user.Id), &ds.Volunteer{
		UserId:          user.Id,
		DemonstrationId: demonstration.Id,
		Polygon:         polygon,
		Location:        location,
		Time:            time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	err = h.wac.SendDemonstrationTemplate(user.Phone, demonstration.Id, user.Id, url.QueryEscape(user.Name), polygon, location)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Volunteer added: %s (%s) polygon %s demonstration %s (%s)", user.Name, user.Phone, polygon, demonstration.Name, demonstration.Id)
	log.Info(msg)
	h.sc.Info(msg)
	return nil
}

func (h *Handle) join(phone string) error {

	user, err := ds.GetUserByPhone(h.dsc, phone)
	if err != nil {
		return err
	}

	return h.Join(user)
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

func getAvailablePolygons(dsc ds.Client, demonstration string) (map[string]string, error) {

	res := getPolygons()
	var volunteers []ds.Volunteer
	err := dsc.GetFilter(ds.KindVolunteer,
		[]ds.FilterField{{Name: "demonstration_id", Operator: "=", Value: demonstration}},
		&volunteers)
	if err != nil {
		return nil, err
	}
	for _, currVolunteer := range volunteers {
		delete(res, currVolunteer.Polygon)
	}

	return res, nil
}

func getPolygons() map[string]string {

	return map[string]string{
		"A1":   "32.072735,34.793577",
		"A2":   "32.072885,34.793057",
		"A3":   "32.072885,34.793057",
		"A4":   "32.07323,34.791968",
		"A5":   "32.073358,34.791222",
		"A6":   "32.073299,34.790208",
		"A7":   "32.073776,34.790851",
		"A8":   "32.074244,34.791028",
		"A9":   "32.074758,34.791253",
		"A10":  "32.074944,34.790765",
		"A11":  "32.074361,34.79045",
		"A12":  "32.074016,34.790332",
		"A13":  "32.073716,34.790257",
		"A14":  "32.072565,34.789823",
		"A15":  "32.071947,34.789538",
		"A16":  "32.072052,34.79008",
		"A17":  "32.072602,34.790322",
		"A18":  "32.075203,34.79121",
		"A19":  "32.075576,34.791255",
		"A20":  "32.075481,34.791486",
		"A21":  "32.071565,34.78979",
		"A22":  "32.071433,34.789254",
		"A23":  "32.072971,34.789724",
		"A24":  "32.07607,34.791718",
		"K1":   "32.07333,34.789284",
		"K2":   "32.07333,34.78879",
		"K3":   "32.07333,34.788393",
		"K4":   "32.07333,34.788002",
		"K5":   "32.07333,34.787701",
		"K6":   "32.073312,34.787385",
		"K7":   "32.07333,34.78709",
		"K8":   "32.073326,34.786784",
		"K9":   "32.07333,34.786521",
		"K10":  "32.073344,34.786178",
		"K11":  "32.073329,34.785806",
		"K12":  "32.073388,34.785462",
		"K13":  "32.073384,34.785098",
		"K14":  "32.073384,34.784786",
		"K15":  "32.073688,34.784781",
		"K16":  "32.073934,34.78477",
		"K17":  "32.073056,34.784797",
		"K18":  "32.072738,34.784807",
		"K18B": "32.072974,34.785312",
		"K19":  "32.073406,34.784502",
		"K20":  "32.073411,34.784212",
		"K21":  "32.073411,34.783944",
		"K22":  "32.073406,34.783649",
		"K23":  "32.073702,34.783676",
		"K24":  "32.073997,34.783676",
		"K25":  "32.073124,34.783631",
		"K26":  "32.073406,34.78332",
		"K27":  "32.073415,34.783019",
		"K28":  "32.073484,34.782488",
		"L3":   "32.072592,34.783522",
		"L4":   "32.072865,34.783607",
		"G1":   "32.07399,34.781832",
		"G2":   "32.073581,34.781837",
		"G3":   "32.073194,34.781821",
		"G4":   "32.072935,34.781794",
		"D1":   "32.073657,34.781397",
	}
}
