package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Register struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (h *Handle) Register(w http.ResponseWriter, r *http.Request) {

	var request Register
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Infof("failed to decode request registration with '%v'", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !validateRegisterRequest(h.dsc, &request) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if string(request.Phone[0]) == "0" {
		request.Phone = fmt.Sprintf("972%s", request.Phone[1:])
	}

	token := uuid.NewString()
	err = h.dsc.Put(ds.KindRegisterRequest, token, &ds.RegisterRequest{
		Phone: request.Phone,
		Name:  request.Name,
		Time:  time.Now().Unix(),
	})
	if err != nil {
		log.Errorf("failed to create user '%+v' in datastore with '%v'", request, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.wac.SendSignupTemplate(request.Phone, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func validateRegisterRequest(dsc ds.Client, request *Register) bool {

	if len(request.Name) > 32 {
		log.Info("invalid register name")
		return false
	}

	if !regexp.MustCompile(`^[0-9]{10}$`).MatchString(request.Phone) {
		log.Infof("invalid phone '%s'", request.Phone)
		return false
	}

	return dsc.GetFilter(ds.KindRegisterRequest, "phone", "=", request.Phone, &[]ds.RegisterRequest{}) == nil
}
