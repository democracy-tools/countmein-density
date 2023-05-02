package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (h *Handle) DeleteUser(w http.ResponseWriter, r *http.Request) {

	userId := mux.Vars(r)["user-id"]
	if !validateToken(userId) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user ds.User
	err := h.dsc.Get(ds.KindUser, userId, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.dsc.Delete(ds.KindUser, userId)
	if err != nil {
		h.sc.Debug(fmt.Sprintf("Failed to delete user %s (%s) %s with %v", user.Name, user.Phone, userId, err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.sc.Info(fmt.Sprintf("User deleted %s (%s) %s", user.Name, user.Phone, userId))
}

func deleteUser(dsc ds.Client, phone string) error {

	user, err := ds.GetUserByPhone(dsc, phone)
	if err != nil {
		return err
	}
	if user == nil {
		err = fmt.Errorf("user with phone '%s' not found", phone)
		log.Info(err.Error())
		return err
	}

	return dsc.Delete(ds.KindUser, user.Id)
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

	err = wac.SendOnboardingTemplate(phone, id)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusCreated
}

func validateUser(dsc ds.Client, phone string, name string) int {

	user, err := ds.GetUserByPhone(dsc, phone)
	if err != nil {
		return http.StatusInternalServerError
	}
	if user != nil {
		log.Infof("user %s (%s) already exist", name, phone)
		return http.StatusBadRequest
	}

	return http.StatusOK
}
