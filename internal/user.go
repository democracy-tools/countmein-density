package internal

import (
	"fmt"
	"net/http"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/gorilla/mux"
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
