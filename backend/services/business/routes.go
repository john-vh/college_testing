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
	userIdParam     = "userId"
)

func (h *BusinessHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /admin/businesses", h.handleErr(h.handleQueryAllBusinesses))
	router.HandleFunc("GET /admin/posts", h.handleErr(h.handleQueryAllPosts))
	router.HandleFunc("POST /admin/businesses/{businessId}/approve", h.handleErr(h.handleApproveBusiness))

	router.HandleFunc("GET /posts", h.handleErr(h.handleGetActivePosts))

	router.HandleFunc("GET /businesses", h.handleErr(h.handleGetBusinesses))
	router.HandleFunc("GET /businesses/{businessId}", h.handleErr(h.handleGetBusiness))
	router.HandleFunc("GET /users/0/businesses", h.handleErr(h.handleGetUserBusinesses))
	router.HandleFunc("GET /users/0/posts", h.handleErr(h.handleGetUserPosts))
	router.HandleFunc("GET /users/0/applications", h.handleErr(h.handleGetUserApplications))
	router.HandleFunc("POST /users/0/businesses", h.handleErr(h.handleRequestBusiness))
	router.HandleFunc("PATCH /businesses/{businessId}", h.handleErr(h.handleUpdateBusiness))

	router.HandleFunc("POST /businesses/{businessId}/posts", h.handleErr(h.handleCreatePost))
	router.HandleFunc("PATCH /businesses/{businessId}/posts/{postId}", h.handleErr(h.handleUpdatePost))
	router.HandleFunc("POST /businesses/{businessId}/posts/{postId}/activate", h.handleErr(h.handleActivatePost))
	router.HandleFunc("POST /businesses/{businessId}/posts/{postId}/deactivate", h.handleErr(h.handleDeactivatePost))
	router.HandleFunc("POST /businesses/{businessId}/posts/{postId}/archive", h.handleErr(h.handleArchivePost))

	router.HandleFunc("POST /businesses/{businessId}/posts/{postId}/apply", h.handleErr(h.handleApplyToPost))
	router.HandleFunc("GET /businesses/{businessId}/posts/{postId}/applications", h.handleErr(h.handleGetPostApplications))
	router.HandleFunc("POST /businesses/{businessId}/posts/{postId}/applications/{userId}/accept", h.handleErr(h.handleAcceptApplication))
	router.HandleFunc("POST /businesses/{businessId}/posts/{postId}/applications/{userId}/reject", h.handleErr(h.handleRejectApplication))
	router.HandleFunc("POST /businesses/{businessId}/posts/{postId}/applications/{userId}/complete", h.handleErr(h.handleCompleteApplication))
	router.HandleFunc("POST /businesses/{businessId}/posts/{postId}/applications/{userId}/incomplete", h.handleErr(h.handleAbandonApplication))
	router.HandleFunc("POST /businesses/{businessId}/posts/{postId}/applications/{userId}/withdraw", h.handleErr(h.handleWithdrawApplication))
}

func (h *BusinessHandler) handleQueryAllBusinesses(w http.ResponseWriter, r *http.Request) error {
	const (
		param_user string = "user"
	)
	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	var userId *uuid.UUID
	if r.URL.Query().Has(param_user) {
		if id, err := uuid.Parse(r.URL.Query().Get(param_user)); err == nil {
			userId = &id
		}
	}

	// TODO: Paramaterize status
	params := models.BusinessQueryParams{
		UserId: userId,
	}

	businesses, err := h.GetBusinesses(r.Context(), session, &params)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(businesses)
	return nil
}

func (h *BusinessHandler) handleQueryAllPosts(w http.ResponseWriter, r *http.Request) error {
	const (
		param_business string = "business"
		param_user     string = "user"
	)
	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	var businessId *uuid.UUID
	if r.URL.Query().Has(param_business) {
		if id, err := uuid.Parse(r.URL.Query().Get(param_business)); err == nil {
			businessId = &id
		}
	}
	var userId *uuid.UUID
	if r.URL.Query().Has(param_user) {
		if id, err := uuid.Parse(r.URL.Query().Get(param_user)); err == nil {
			userId = &id
		}
	}

	// TODO: Paramaterize status
	params := models.PostQueryParams{
		BusinessId: businessId,
		UserId:     userId,
	}

	posts, err := h.GetPosts(r.Context(), session, &params)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
	return nil
}

func (h *BusinessHandler) handleGetActivePosts(w http.ResponseWriter, r *http.Request) error {
	const (
		param_business string = "business"
	)
	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	var businessId *uuid.UUID
	if r.URL.Query().Has(param_business) {
		if id, err := uuid.Parse(r.URL.Query().Get(param_business)); err == nil {
			businessId = &id
		}
	}

	status := models.POST_STATUS_ACTIVE
	params := models.PostQueryParams{
		Status:     &status,
		BusinessId: businessId,
	}

	posts, err := h.GetPosts(r.Context(), session, &params)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
	return nil
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

func (h *BusinessHandler) handleGetBusinesses(w http.ResponseWriter, r *http.Request) error {
	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	status := models.BUSINESS_STATUS_ACTIVE

	params := models.BusinessQueryParams{
		Status: &status,
	}

	businesses, err := h.GetBusinesses(r.Context(), session, &params)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(businesses)
	return nil
}

func (h *BusinessHandler) handleGetBusiness(w http.ResponseWriter, r *http.Request) error {
	businessId, err := uuid.Parse(r.PathValue(businessIdParam))
	if err != nil {
		return services.NewNotFoundServiceError(err)
	}

	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	business, err := h.GetBusinessForId(r.Context(), session, &businessId)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(business)
	return nil
}

func (h *BusinessHandler) handleGetUserBusinesses(w http.ResponseWriter, r *http.Request) error {
	const (
		param_status string = "status"
	)

	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	params := models.BusinessQueryParams{
		UserId: session.GetUserId(),
	}

	businesses, err := h.GetBusinesses(r.Context(), session, &params)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(businesses)
	return nil
}

func (h *BusinessHandler) handleGetUserPosts(w http.ResponseWriter, r *http.Request) error {
	const (
		param_business = "business"
	)
	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	var businessId *uuid.UUID
	if r.URL.Query().Has(param_business) {
		if id, err := uuid.Parse(r.URL.Query().Get(param_business)); err == nil {
			businessId = &id
		}
	}

	params := models.PostQueryParams{
		UserId:     session.GetUserId(),
		BusinessId: businessId,
	}

	posts, err := h.GetPosts(r.Context(), session, &params)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
	return nil
}

func (h *BusinessHandler) handleGetUserApplications(w http.ResponseWriter, r *http.Request) error {
	session, err := h.sessions.GetSession(r)
	if err != nil {
		return err
	}

	params := models.UserApplicationQueryParams{
		UserId: session.GetUserId(),
	}

	applications, err := h.GetUserApplications(r.Context(), session, &params)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)
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

func (h *BusinessHandler) handleSetPostStatus(status models.PostStatus) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
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

		err = h.SetPostStatus(r.Context(), session, &businessId, postId, status)
		if err != nil {
			return err
		}

		return nil
	}
}

func (h *BusinessHandler) handleActivatePost(w http.ResponseWriter, r *http.Request) error {
	return h.handleSetPostStatus(models.POST_STATUS_ACTIVE)(w, r)
}

func (h *BusinessHandler) handleDeactivatePost(w http.ResponseWriter, r *http.Request) error {
	return h.handleSetPostStatus(models.POST_STATUS_DISABLED)(w, r)
}

func (h *BusinessHandler) handleArchivePost(w http.ResponseWriter, r *http.Request) error {
	return h.handleSetPostStatus(models.POST_STATUS_ARCHIVED)(w, r)
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

func (h *BusinessHandler) handleSetApplicationStatus(status models.ApplicationStatus) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		businessId, err := uuid.Parse(r.PathValue(businessIdParam))
		if err != nil {
			return services.NewNotFoundServiceError(err)
		}

		postId, err := strconv.Atoi(r.PathValue(postIdParam))
		if err != nil {
			return services.NewNotFoundServiceError(err)
		}

		userId, err := uuid.Parse(r.PathValue(userIdParam))
		if err != nil {
			return services.NewNotFoundServiceError(err)
		}

		session, err := h.sessions.GetSession(r)
		if err != nil {
			return err
		}

		err = h.SetApplicationStatus(r.Context(), session, &businessId, postId, &userId, status)
		if err != nil {
			return err
		}

		return nil
	}
}

func (h *BusinessHandler) handleAcceptApplication(w http.ResponseWriter, r *http.Request) error {
	return h.handleSetApplicationStatus(models.APPLICATION_STATUS_ACCEPTED)(w, r)
}

func (h *BusinessHandler) handleRejectApplication(w http.ResponseWriter, r *http.Request) error {
	return h.handleSetApplicationStatus(models.APPLICATION_STATUS_REJECTED)(w, r)
}

func (h *BusinessHandler) handleCompleteApplication(w http.ResponseWriter, r *http.Request) error {
	return h.handleSetApplicationStatus(models.APPLICATION_STATUS_COMPLETED)(w, r)
}

func (h *BusinessHandler) handleAbandonApplication(w http.ResponseWriter, r *http.Request) error {
	return h.handleSetApplicationStatus(models.APPLICATION_STATUS_CANCELLED)(w, r)
}

func (h *BusinessHandler) handleWithdrawApplication(w http.ResponseWriter, r *http.Request) error {
	return h.handleSetApplicationStatus(models.APPLICATION_STATUS_WITHDRAWN)(w, r)
}
