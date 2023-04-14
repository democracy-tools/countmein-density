package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Observation struct {
	Time          int64   `json:"time" datastore:"time"`
	User          string  `json:"user_id" datastore:"user_id"`
	Demonstration string  `json:"demonstration" datastore:"demonstration"`
	Polygon       string  `json:"polygon" datastore:"polygon"`
	Density       float32 `json:"density" datastore:"density"`
	Latitude      float32 `json:"latitude" datastore:"latitude"`
	Longitude     float32 `json:"longitude" datastore:"longitude"`
}

func (h *Handle) CreateObservation(w http.ResponseWriter, r *http.Request) {

	var observation Observation
	err := json.NewDecoder(r.Body).Decode(&observation)
	if err != nil {
		log.Infof("failed to decode request observation with '%v'", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	observation.Time = time.Now().Unix()
	if !validateObservation(&observation) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.dsc.Put(ds.KindObservation, uuid.NewString(), &observation)
	if err != nil {
		log.Errorf("failed to insert new observation '%+v' into datastore with '%v'", observation, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handle) GetObservations(w http.ResponseWriter, r *http.Request) {

	var observations []Observation
	err := h.dsc.GetFilter(ds.KindObservation, "time", ">", time.Now().Add(time.Minute*(-17)).Unix(), &observations)
	if err != nil {
		log.Errorf("failed to get observations by time with '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	observations = latestObservations(observations)
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"observations": observations})
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(getObservationAsText(observations))
	}
}

func latestObservations(observations []Observation) []Observation {

	userToLastObservation := make(map[string]Observation)
	for _, currObservation := range observations {
		latest, ok := userToLastObservation[currObservation.Polygon]
		if !ok || latest.Time < currObservation.Time {
			userToLastObservation[currObservation.Polygon] = currObservation
		}
	}

	return toObservationSlice(userToLastObservation)
}

func toObservationSlice(userToLastObservation map[string]Observation) []Observation {

	var res []Observation
	for _, currObservation := range userToLastObservation {
		res = append(res, currObservation)
	}

	return res
}

func getObservationAsText(observations []Observation) []byte {

	sort.Slice(observations, func(i, j int) bool {
		return observations[i].Polygon < observations[j].Polygon
	})

	var buf bytes.Buffer
	for _, currObservation := range observations {
		buf.WriteString(fmt.Sprintf("%s: %.1f\n", currObservation.Polygon, currObservation.Density))
	}

	res := buf.Bytes()
	if len(res) == 0 {
		return []byte("No observation found")
	}
	return res
}

func validateObservation(observation *Observation) bool {

	now := time.Now()
	if observation.Time < now.Add(time.Hour*(-2)).Unix() || observation.Time > now.Add(time.Hour*2).Unix() {
		log.Infof("invalid observation time '%d' user '%s'", observation.Time, observation.User)
		return false
	}

	if len(observation.User) < 2 {
		log.Infof("invalid observation user '%s'", observation.User)
		return false
	}

	match, err := regexp.MatchString("^[A-Z]+[1-9][0-9]{0,2}[A-Z]?$", observation.Polygon)
	if err != nil {
		log.Errorf("failed to validate polygon '%s' using regexp with '%v'", observation.Polygon, err)
		return false
	}
	if !match {
		log.Infof("invalid observation polygon '%s' user '%s'", observation.Polygon, observation.User)
		return false
	}

	if observation.Density < 0 && observation.Density > 9 {
		log.Infof("invalid observation density '%f' user '%s'", observation.Density, observation.User)
		return false
	}

	if observation.Latitude < 0 && observation.Latitude > 40 {
		log.Infof("invalid observation latitude '%f' user '%s'", observation.Latitude, observation.User)
		return false
	}

	if observation.Longitude < 0 && observation.Longitude > 40 {
		log.Infof("invalid observation longitude '%f' user '%s'", observation.Longitude, observation.User)
		return false
	}

	return true
}
