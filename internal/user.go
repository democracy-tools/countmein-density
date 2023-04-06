package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (h *Handle) CreateUser(w http.ResponseWriter, r *http.Request) {

	request, code := getRegisterRequest(h.dsc, r)
	if code != http.StatusOK {
		w.WriteHeader(code)
		return
	}

	id := uuid.NewString()
	err := h.dsc.Put(ds.KindUser, id, &ds.User{
		Id:         id,
		Phone:      request.Phone,
		Name:       request.Name,
		Preference: request.Preference,
		Time:       time.Now().Unix(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.wac.Send(request.Phone, fmt.Sprintf("האימות עבר בהצלחה, נשלח לך הודעה לפני ההפגנה עם פרטים.\nבנתיים את/ה מוזמן/ת להצטרף גם לקבוצת הוואטאפ שלנו\n%s", WhatAppGroupLink))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
