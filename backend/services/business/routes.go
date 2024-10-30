package business

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services"
)

const (
	businessIdParam = "businessId"
	postIdParam     = "postId"
)

func (h *BusinessHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /businesses", h.handleErr(h.handleRequestBusiness))
	router.HandleFunc("PATCH /businesses/{businessId}", h.handleErr(h.handleUpdateBusiness))
	router.HandleFunc("POST /businesses/{businessId}/approve", h.handleErr(h.handleApproveBusiness))

	router.HandleFunc("POST /businesses/{businessId}/posts", h.handleErr(h.handleCreatePost))
	router.HandleFunc("PATCH /businesses/{businessId}/posts/{postId}", h.handleErr(h.handleUpdatePost))

	router.HandleFunc("POST /businesses/{businessId}/posts/{postId}/apply", h.handleErr(h.handleApplyToPost))
	router.HandleFunc("GET /businesses/{businessId}/posts/{postId}/applications", h.handleErr(h.handleGetPostApplications))
}

func (h *BusinessHandler) handleRequestBusiness(w http.ResponseWriter, r *http.Request) error {
	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}
	userId := session.GetUserId()
	if userId == nil {
		return services.NewUnauthenticatedServiceError(err)
	}

	data := models.BusinessCreate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		return err
	}
	business, err := h.RequestBusiness(r.Context(), session, userId, &data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(business)
	return nil
}

func (h *BusinessHandler) handleUpdateBusiness(_ http.ResponseWriter, r *http.Request) error {
	businessId, err := uuid.Parse(r.PathValue(businessIdParam))
	if err != nil {
		return services.NewNotFoundServiceError(err)
	}

	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	data := models.BusinessUpdate{}
	err = models.ReadRequestJson(r, &data)
	if err != nil {
		return err
	}

	return h.UpdateBusiness(r.Context(), session, &businessId, &data)
}

func (h *BusinessHandler) handleApproveBusiness(_ http.ResponseWriter, r *http.Request) error {
	businessId, err := uuid.Parse(r.PathValue(businessIdParam))
	if err != nil {
		return services.NewNotFoundServiceError(err)
	}

	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	return h.ApproveBusiness(r.Context(), session, &businessId)
}

func (h *BusinessHandler) handleCreatePost(w http.ResponseWriter, r *http.Request) error {
	businessId, err := uuid.Parse(r.PathValue(businessIdParam))
	if err != nil {
		return services.NewNotFoundServiceError(err)
	}

	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	data := models.PostCreate{}
	if err := models.ReadRequestJson(r, &data); err != nil {
		return err
	}

	post, err := h.CreatePost(r.Context(), session, &businessId, &data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
	return nil
}

func (h *BusinessHandler) handleUpdatePost(w http.ResponseWriter, r *http.Request) error {
	businessId, err := uuid.Parse(r.PathValue(businessIdParam))
	if err != nil {
		return services.NewNotFoundServiceError(err)
	}

	postId, err := strconv.Atoi(r.PathValue(postIdParam))
	if err != nil {
		return services.NewNotFoundServiceError(err)
	}

	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	data := models.PostUpdate{}
	if err := models.ReadRequestJson(r, &data); err != nil {
		return err
	}

	err = h.UpdatePost(r.Context(), session, &businessId, postId, &data)
	if err != nil {
		return err
	}

	return nil
}

func (h *BusinessHandler) handleApplyToPost(w http.ResponseWriter, r *http.Request) error {
	businessId, err := uuid.Parse(r.PathValue(businessIdParam))
	if err != nil {
		return services.NewNotFoundServiceError(err)
	}

	postId, err := strconv.Atoi(r.PathValue(postIdParam))
	if err != nil {
		return services.NewNotFoundServiceError(err)
	}

	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}
	userId := session.GetUserId()
	if userId == nil {
		return services.NewUnauthenticatedServiceError(nil)
	}

	err = h.CreateApplication(r.Context(), session, &businessId, postId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (h *BusinessHandler) handleGetPostApplications(w http.ResponseWriter, r *http.Request) error {
	businessId, err := uuid.Parse(r.PathValue(businessIdParam))
	if err != nil {
		return services.NewNotFoundServiceError(err)
	}

	postId, err := strconv.Atoi(r.PathValue(postIdParam))
	if err != nil {
		return services.NewNotFoundServiceError(err)
	}

	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}
	applications, err := h.GetPostApplications(r.Context(), session, &businessId, postId)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)
	return nil
}
