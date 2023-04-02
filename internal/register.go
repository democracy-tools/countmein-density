package internal

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	log "github.com/sirupsen/logrus"
)

type Register struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *Handle) Register(w http.ResponseWriter, r *http.Request) {

	var request Register
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Infof("failed to decode request registration with '%v'", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !validateRegisterRequest(&request) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.client.Put(ds.KindUser, request.Email, &ds.User{
		Email:    request.Email,
		Name:     request.Name,
		Time:     time.Now().Unix(),
		Verified: -1,
	}); err != nil {
		log.Errorf("failed to create user '%+v' in datastore with '%v'", request, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: send email

	w.WriteHeader(http.StatusCreated)
}

func validateRegisterRequest(request *Register) bool {

	if len(request.Name) > 32 {
		log.Info("invalid register name")
		return false
	}

	_, err := mail.ParseAddress(request.Email)
	if err != nil {
		log.Infof("invalid register email with '%v'", err)
		return false
	}

	return true
}
