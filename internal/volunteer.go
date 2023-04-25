package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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

	userId, newPolygon := mux.Vars(r)["user-id"], mux.Vars(r)["polygon"]
	if !validateToken(userId) || !validatePolygon(newPolygon) {
		logrus.Infof("[ChangePolygon] invalid polygon '%s' user '%s'", newPolygon, userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	demonstration, err := ds.GetKaplanDemonstration(h.dsc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	availablePolygons, err := getAvailablePolygons(h.dsc, demonstration.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, ok := availablePolygons[newPolygon]; !ok {
		w.Write([]byte(fmt.Sprintf("Polygon %s is not available", newPolygon)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var volunteer ds.Volunteer
	err = h.dsc.Get(ds.KindVolunteer, ds.GetVolunteerId(demonstration.Id, userId), &volunteer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	oldPolygon := volunteer.Polygon
	volunteer.Polygon = newPolygon
	err = h.dsc.Put(ds.KindVolunteer, ds.GetVolunteerId(demonstration.Id, userId), &volunteer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var user ds.User
	err = h.dsc.Get(ds.KindUser, userId, &user)
	if err != nil {
		user = ds.User{
			Id:    userId,
			Phone: "na",
			Name:  err.Error(),
		}
	}
	msg := fmt.Sprintf("Volunteer %s (%s) changed polygon from %s to %s demonstration %s",
		user.Name, user.Phone, oldPolygon, newPolygon, demonstration.Id)
	log.Info(msg)
	h.sc.Info(msg)
}
