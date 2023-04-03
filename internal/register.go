package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/email"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Register struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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

	token := uuid.NewString()
	err = h.client.Put(ds.KindRegisterRequest, token, &ds.RegisterRequest{
		Name:  request.Name,
		Email: request.Email,
		Time:  time.Now().Unix(),
	})
	if err != nil {
		log.Errorf("failed to create user '%+v' in datastore with '%v'", request, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = email.GetInstance().Send(request.Email,
		"internal/email/verify.template",
		"CountMeIn verify",
		struct{ Link string }{Link: fmt.Sprintf("https://aaa.com?token=%s", token)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
