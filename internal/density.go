package internal

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Observation struct {
	Time    int64
	Polygon string
	Density int
}

func CreateObservation(w http.ResponseWriter, r *http.Request) {

	var o *Observation
	err := json.NewDecoder(r.Body).Decode(o)
	if err != nil {
		log.Infof("failed to decode request observation with '%v'", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func GetObservations(w http.ResponseWriter, r *http.Request) {
}
