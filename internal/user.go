package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatapp"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (h *Handle) CreateUser(w http.ResponseWriter, r *http.Request) {

	request, code := getRegisterRequest(h.dsc, r)
	if code != http.StatusOK {
		w.WriteHeader(code)
		return
	}

	w.WriteHeader(createUser(h.dsc, h.wac, request.Phone, request.Name, request.Preference))
}

func createUser(dsc ds.Client, wac whatsapp.Client, phone string, name string, preference string) int {

	code := validateUser(dsc, phone, name)
	if code != http.StatusOK {
		return code
	}

	id := uuid.NewString()
	err := dsc.Put(ds.KindUser, id, &ds.User{
		Id:         id,
		Phone:      phone,
		Name:       name,
		Preference: preference,
		Time:       time.Now().Unix(),
	})
	if err != nil {
		return http.StatusInternalServerError
	}

	err = wac.SendVerifyTemplate(phone)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusCreated
}

func validateUser(dsc ds.Client, phone string, name string) int {

	var users []ds.User
	err := dsc.GetFilter(ds.KindUser, []ds.FilterField{{Name: "phone", Operator: "=", Value: phone}}, &users)
	if err != nil {
		log.Errorf("failed to get user %s (%s) with '%v'", name, phone, err)
		return http.StatusInternalServerError
	}
	if len(users) > 0 {
		log.Infof("user %s (%s) already exist", name, phone)
		return http.StatusBadRequest
	}

	return http.StatusOK
}

func getRegisterRequest(client ds.Client, r *http.Request) (*ds.RegisterRequest, int) {

	var request struct {
		Token string `json:"token"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Infof("failed to decode create user request with '%v'", err)
		return nil, http.StatusBadRequest
	}
	if !validateToken(request.Token) {
		return nil, http.StatusBadRequest
	}

	var res ds.RegisterRequest
	err = client.Get(ds.KindRegisterRequest, request.Token, &res)
	if err != nil {
		if ds.IsNoSuchEntityError(err) {
			log.Infof("'%s' with token '%s' not found", ds.KindRegisterRequest, request.Token)
			return nil, http.StatusBadRequest
		}
		log.Errorf("failed to get '%s' from datastore with '%v'", ds.KindRegisterRequest, err)
		return nil, http.StatusInternalServerError
	}

	return &res, http.StatusOK
}

func validateToken(token string) bool {

	count := len(token)
	if count < 30 || count > 40 {
		log.Infof("invalid token length '%d'", count)
		return false
	}

	return true
}
