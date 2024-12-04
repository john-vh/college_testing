package user

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/models"
)

func (h *UserHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /users", h.handleErr(h.handleGetUsers))
	router.HandleFunc("PATCH /users/0", h.handleErr(h.handleUpdateUser))
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

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}
	data := models.UserUpdate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		return err
	}

	if err := h.UpdateUser(r.Context(), session, session.GetUserId(), &data); err != nil {
		return err
	}

	return nil
}
