package business

import (
	"context"
	"errors"
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
	userId := session.GetUserId()
	if userId == nil {
		return nil, services.NewUnauthenticatedServiceError(nil)
	}

	err := models.ValidateData(data)
	if err != nil {
		return nil, err
	}

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.Business, error) {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return nil, services.NewUnauthenticatedServiceError(err)
		}

		err = AuthorizeBusinessAction(user, BUSINESS_ACTION_CREATE, nil, nil)
		if err != nil {
			return nil, err
		}

		business, dberr := pq.CreateBusiness(ctx, ownerId, data)
		if dberr != nil {
			if errors.Is(dberr, db.ErrUnique) {
				return nil, services.NewDataConflictServiceError(err, "Business must be unique")
			}
			return nil, err
		}

		return business, nil
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
		if err := AuthorizeBusinessAction(user, BUSINESS_ACTION_READ, nil, params); err != nil {
			return nil, err
		}
		return pq.GetBusinesses(ctx, params)
	})
}

func (h *BusinessHandler) UpdateBusiness(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, data *models.BusinessUpdate) error {
	userId := session.GetUserId()
	if userId == nil {
		return services.NewUnauthenticatedServiceError(nil)
	}

	err := models.ValidateData(data)
	if err != nil {
		return err
	}

	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return services.NewUnauthenticatedServiceError(err)
		}
		business, err := pq.GetBusinessForId(ctx, businessId)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return services.NewNotFoundServiceError(err)
			}
			return err
		}
		if err := AuthorizeBusinessAction(user, BUSINESS_ACTION_UPDATE, business, nil); err != nil {
			return err
		}
		err = pq.UpdateBusiness(ctx, businessId, data)
		if err != nil {
			if errors.Is(err, db.ErrUnique) {
				return services.NewDataConflictServiceError(err, "Business must be unique")
			}
			return err
		}

		return nil
	})

}

func (h *BusinessHandler) ApproveBusiness(ctx context.Context, session *sessions.Session, businessId *uuid.UUID) error {
	userId := session.GetUserId()
	if userId == nil {
		return services.NewUnauthenticatedServiceError(nil)
	}

	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return services.NewNotFoundServiceError(err)
			}
			return err
		}
		business, err := pq.GetBusinessForId(ctx, businessId)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return services.NewNotFoundServiceError(err)
			}
			return err
		}
		if err := AuthorizeBusinessAction(user, BUSINESS_ACTION_APPROVE, business, nil); err != nil {
			return err
		}
		err = pq.SetBusinessStatus(ctx, businessId, models.BUSINESS_STATUS_ACTIVE)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return services.NewNotFoundServiceError(err)
			}
			return err
		}
		return nil
	})
}

type BusinessAction string

const (
	BUSINESS_ACTION_CREATE  BusinessAction = "business:create"
	BUSINESS_ACTION_UPDATE  BusinessAction = "business:update"
	BUSINESS_ACTION_APPROVE BusinessAction = "business:approve"
	BUSINESS_ACTION_READ    BusinessAction = "business:read"
)

func AuthorizeBusinessAction(user *models.User, action BusinessAction, data *models.Business, query *models.BusinessQueryParams) error {
	if user == nil {
		return services.NewUnauthenticatedServiceError(nil)
	}

	for _, role := range user.Roles {
		switch role {
		case models.USER_ROLE_ADMIN:
			switch action {
			case BUSINESS_ACTION_CREATE:
				return nil
			case BUSINESS_ACTION_UPDATE:
				return nil
			case BUSINESS_ACTION_APPROVE:
				return nil
			case BUSINESS_ACTION_READ:
				return nil
			}
		case models.USER_ROLE_USER:
			switch action {
			case BUSINESS_ACTION_CREATE:
				return nil
			case BUSINESS_ACTION_UPDATE:
				if data != nil && data.UserId == user.Id {
					return nil
				}
			case BUSINESS_ACTION_READ:
				if query != nil &&
					((query.UserId != nil && *query.UserId == user.Id) ||
						(query.Status != nil && *query.Status == models.BUSINESS_STATUS_ACTIVE)) {
					return nil
				}
			}
		}
	}

	return services.NewUnauthorizedServiceError(nil)
}
