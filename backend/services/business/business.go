package business

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/filestore"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services"
	"github.com/john-vh/college_testing/backend/services/notifications"
	"github.com/john-vh/college_testing/backend/services/sessions"
	"github.com/john-vh/college_testing/backend/services/user"
)

type BusinessHandler struct {
	logger                    *slog.Logger
	sessions                  *sessions.SessionsHandler
	users                     *user.UserHandler
	store                     *db.PgxStore
	filestore                 filestore.FileStore
	notifications             *notifications.NotificationsService
	notificationsTemplatesDir string
	frontendURL               string
	handleErr                 services.ServicesHTTPErrorHandler
}

func NewBusinessHandler(
	logger *slog.Logger,
	sessions *sessions.SessionsHandler,
	users *user.UserHandler,
	store *db.PgxStore,
	filestore filestore.FileStore,
	notifications *notifications.NotificationsService,
	notificationsTemplatesDir string,
	frontendURL string,
	errHandler services.ServicesHTTPErrorHandler,
) *BusinessHandler {
	return &BusinessHandler{
		logger:                    logger,
		sessions:                  sessions,
		users:                     users,
		store:                     store,
		filestore:                 filestore,
		notifications:             notifications,
		notificationsTemplatesDir: notificationsTemplatesDir,
		frontendURL:               frontendURL,
		handleErr:                 errHandler,
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

	business, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.Business, error) {
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
	if err != nil {
		return nil, err
	}

	// Asyncronously send notifications
	go h.sendBusinessRequestedNotifications(business)

	return business, nil
}

func (h *BusinessHandler) setBusinessImage(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, filename string, f io.ReadSeeker) (err error) {
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
		err = AuthorizeBusinessAction(user, BUSINESS_ACTION_UPDATE, business, nil)
		if err != nil {
			return err
		}
		ext := filepath.Ext(filename)
		key := fmt.Sprintf("%v%v", businessId.String(), ext)
		prevKey := ""
		if business.LogoUrl != nil {
			prevKey = h.filestore.GetKey(*business.LogoUrl)
		}
		url := h.filestore.GetURI(key)
		err = pq.SetBusinessLogo(ctx, businessId, url)
		if err != nil {
			return err
		}
		err = h.filestore.UploadObject(key, f)
		if err != nil {
			h.logger.Warn("Failed to upload image for business", "err", err)
			return err
		}

		if prevKey != key && prevKey != "" {
			err = h.filestore.DeleteObject(prevKey)
			if err != nil {
				h.logger.Warn("Failed to delete old business logo", "err", err)
				return err
			}
		}

		return nil
	})
}

func (h *BusinessHandler) sendBusinessRequestedNotifications(b *models.Business) error {
	var admins []models.User
	var owner *models.User

	err := db.WithTx(context.Background(), h.store, func(pq *db.PgxQueries) error {
		var err error
		owner, err = pq.GetUserForId(context.Background(), &b.UserId)
		if err != nil {
			return err
		}

		status := models.USER_STATUS_ACTIVE
		role := models.USER_ROLE_ADMIN
		admins, err = pq.QueryUsers(context.Background(), &models.UserQueryParams{Status: &status, Role: &role})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		h.logger.Warn("Failed to enqueue business requested notifications", "err", err)
		return err
	}

	for _, admin := range admins {
		err := h.notifications.Enqueue(context.Background(), h.NewBusinessRequestedAdminNotification(&admin, owner, b))
		if err != nil {
			h.logger.Warn("Failed to enqueue business requested admin notification", "err", err)
		}
	}

	return nil
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

func (h *BusinessHandler) GetBusinessForId(ctx context.Context, session *sessions.Session, id *uuid.UUID) (*models.Business, error) {
	userId := session.GetUserId()
	if userId == nil {
		return nil, services.NewUnauthenticatedServiceError(nil)
	}

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.Business, error) {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return nil, services.NewUnauthenticatedServiceError(err)
		}
		business, err := pq.GetBusinessForId(ctx, id)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return nil, services.NewNotFoundServiceError(err)
			}
			return nil, err
		}
		if err := AuthorizeBusinessAction(user, BUSINESS_ACTION_READ, business, nil); err != nil {
			return nil, err
		}
		return business, nil
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
				if (query != nil &&
					((query.UserId != nil && *query.UserId == user.Id) ||
						(query.Status != nil && *query.Status == models.BUSINESS_STATUS_ACTIVE))) ||
					(data != nil && data.Status == models.BUSINESS_STATUS_ACTIVE) {
					return nil
				}
			}
		}
	}

	return services.NewUnauthorizedServiceError(nil)
}
