package user

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (h *UserHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /users", h.handleErr(h.handleGetUsers))
}

func (h *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	const idParam = "id"
	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	var userId *uuid.UUID
	if r.URL.Query().Has(idParam) {
		if in, err := uuid.Parse(r.URL.Query().Get(idParam)); err == nil {
			userId = &in
		}
	} else {
		userId = session.GetUserId()
	}

	user, err := h.GetUserById(r.Context(), session, userId)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	return nil
}
