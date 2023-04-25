package internal

import (
	"encoding/json"
	"net/http"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (h *Handle) GetVolunteer(w http.ResponseWriter, r *http.Request) {

	userId := mux.Vars(r)["user-id"]
	if !validateToken(userId) {
		logrus.Infof("[GetVolunteer] invalid user '%s'", userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	demonstration, err := ds.GetKaplanDemonstration(h.dsc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var volunteer ds.Volunteer
	err = h.dsc.Get(ds.KindVolunteer, ds.GetVolunteerId(demonstration.Id, userId), &volunteer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(volunteer)
	if err != nil {
		logrus.Errorf("failed to encode volunteer '%+v' with '%v'", volunteer, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handle) ChangePolygon(w http.ResponseWriter, r *http.Request) {

	// userId, polygon := mux.Vars(r)["user-id"], mux.Vars(r)["polygon"]
	// if !validatePolygon(polygon) {
	// 	logrus.Infof("[ChangePolygon] invalid polygon '%s' user '%s'", polygon, userId)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// demonstration, err := ds.GetKaplanDemonstration(h.dsc)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// availablePolygons, err := getAvailablePolygons(h.dsc, demonstration.Id)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// if _, ok := availablePolygons[polygon]; !ok {
	// 	w.Write([]byte(fmt.Sprintf("Polygon %s is not available", polygon)))
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// var volunteer ds.Volunteer
	// err = h.dsc.Get(ds.KindVolunteer, ds.GetVolunteerId(demonstration.Id, userId), &volunteer)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// assignedPolygon := volunteer.Polygon
	// volunteer.Polygon = polygon
	// err = h.dsc.Put(ds.KindVolunteer, ds.GetVolunteerId(demonstration.Id, userId), &volunteer)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// err = h.wac.SendDemonstrationTemplate(user.Phone, demonstration.Id, user.Id, url.QueryEscape(user.Name), polygon, location)
	// if err != nil {
	// 	return err
	// }

	// msg := fmt.Sprintf("Volunteer %s (%s) changed polygon from %s to %s demonstration %s", user.Name, user.Phone, polygon, demonstration.Name, demonstration.Id)
	// log.Info(msg)
	// h.sc.Info(msg)
}
