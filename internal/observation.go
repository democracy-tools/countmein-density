package internal

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Observation struct {
	Time    int64   `json:"time" datastore:"time"`
	User    string  `json:"user" datastore:"user"`
	Polygon string  `json:"polygon" datastore:"polygon"`
	Density float32 `json:"density" datastore:"density"`
}

type Handle struct{ client ds.Client }

func NewHandle(client ds.Client) *Handle {

	return &Handle{client: client}
}

func (h *Handle) CreateObservation(w http.ResponseWriter, r *http.Request) {

	var observation Observation
	err := json.NewDecoder(r.Body).Decode(&observation)
	if err != nil {
		log.Infof("failed to decode request observation with '%v'", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !validateObservation(&observation) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.client.Put(ds.KindObservation, uuid.NewString(), &observation)
	if err != nil {
		log.Errorf("failed to insert new observation '%+v' into datastore with '%v'", observation, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handle) GetObservations(w http.ResponseWriter, r *http.Request) {

	var observations []Observation
	err := h.client.GetByTime(ds.KindObservation, time.Now().Add(time.Hour*(-3)).Unix(), &observations)
	if err != nil {
		log.Errorf("failed to get observations by time with '%v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"observations": observations})
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

	match, err := regexp.MatchString("^[A-Z]+[1-9][0-9]{0,2}$", observation.Polygon)
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

	return true
}
