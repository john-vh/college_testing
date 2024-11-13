package business

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services"
	"github.com/john-vh/college_testing/backend/services/notifications"
	"github.com/john-vh/college_testing/backend/services/sessions"
	"github.com/john-vh/college_testing/backend/services/user"
)

type BusinessHandler struct {
	logger        *slog.Logger
	sessions      *sessions.SessionsHandler
	notifications *notifications.MailClient
	users         *user.UserHandler
	store         *db.PgxStore
	handleErr     services.ServicesHTTPErrorHandler
}

func NewBusinessHandler(
	logger *slog.Logger,
	errHandler services.ServicesHTTPErrorHandler,
	sessions *sessions.SessionsHandler,
	users *user.UserHandler,
	notifications *notifications.MailClient,
	store *db.PgxStore,
) *BusinessHandler {
	return &BusinessHandler{
		logger:        logger,
		sessions:      sessions,
		notifications: notifications,
		store:         store,
		handleErr:     errHandler,
		users:         users,
	}
}

func (h *BusinessHandler) RequestBusiness(ctx context.Context, session *sessions.Session, ownerId *uuid.UUID, data *models.BusinessCreate) (*models.Business, error) {
	err := models.ValidateData(data)
	if err != nil {
		return nil, err
	}

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.Business, error) {
		_, err := h.users.AuthorizeModifyUser(ctx, pq, session, ownerId)
		if err != nil {
			return nil, err
		}

		return pq.CreateBusiness(ctx, ownerId, data)
	})
}

func (h *BusinessHandler) GetBusinesses(ctx context.Context, session *sessions.Session, params *models.BusinessQueryParams) ([]models.Business, error) {
	userId := session.GetUserId()
	if userId == nil {
		return nil, services.NewUnauthenticatedServiceError(nil)
	}

	if params == nil {
		params = &models.BusinessQueryParams{}
	}

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.Business, error) {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return nil, services.NewUnauthenticatedServiceError(err)
		}

		if (params.Status == nil || *params.Status != models.BUSINESS_STATUS_ACTIVE) &&
			!(user.HasRole(models.USER_ROLE_ADMIN) ||
				(params.UserId != nil && *params.UserId == *userId)) {
			return nil, services.NewUnauthorizedServiceError(fmt.Errorf("User attempted to retrieve non-active businesses"))
		}

		return pq.GetBusinesses(ctx, params)
	})
}

func (h *BusinessHandler) UpdateBusiness(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, data *models.BusinessUpdate) error {
	err := models.ValidateData(data)
	if err != nil {
		return err
	}

	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		_, err := h.AuthorizeModifyBusiness(ctx, pq, session, businessId)
		if err != nil {
			return err
		}
		return pq.UpdateBusiness(ctx, businessId, data)
	})

}

func (h *BusinessHandler) ApproveBusiness(ctx context.Context, session *sessions.Session, businessId *uuid.UUID) error {
	sUserId := session.GetUserId()
	if sUserId == nil {
		return services.NewUnauthenticatedServiceError(nil)
	}

	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		user, err := pq.GetUserForId(ctx, sUserId)
		if err != nil {
			return err
		}

		if !user.HasRole(models.USER_ROLE_ADMIN) {
			return services.NewUnauthorizedServiceError(nil)
		}

		return pq.SetBusinessStatus(ctx, businessId, models.BUSINESS_STATUS_ACTIVE)
	})
}

func (h *BusinessHandler) AuthorizeModifyBusiness(ctx context.Context, pq *db.PgxQueries, session *sessions.Session, businessId *uuid.UUID) (*models.User, error) {
	userId := session.GetUserId()
	if userId == nil {
		return nil, services.NewUnauthorizedServiceError(nil)
	}

	owner, err := pq.GetBusinessOwner(ctx, businessId)
	if err != nil {
		// TODO: Handle db error
		return nil, err
	}

	if owner.Id == *userId {
		return owner, nil
	}

	user, err := pq.GetUserForId(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user.HasRole(models.USER_ROLE_ADMIN) {
		return user, nil
	}

	return nil, services.NewUnauthorizedServiceError(nil)
}
