package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services"
	"github.com/john-vh/college_testing/backend/services/sessions"
)

type UserHandler struct {
	logger    *slog.Logger
	sessions  *sessions.SessionsHandler
	store     *db.PgxStore
	handleErr services.ServicesHTTPErrorHandler
}

func NewUserHandler(
	logger *slog.Logger,
	errHandler services.ServicesHTTPErrorHandler,
	sessions *sessions.SessionsHandler,
	store *db.PgxStore,
) *UserHandler {
	return &UserHandler{
		logger:    logger,
		sessions:  sessions,
		store:     store,
		handleErr: errHandler,
	}
}

func (h *UserHandler) GetUserById(ctx context.Context, session *sessions.Session, id *uuid.UUID) (*models.User, error) {
	user, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.User, error) {
		return h.AuthorizeModifyUser(ctx, pq, session, id)
	})
	if err != nil {
		h.logger.Debug("Error getting user", "err", err)
		return nil, err
	}

	return user, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, session *sessions.Session, id *uuid.UUID, data *models.UserUpdate) error {
	if err := models.ValidateData(data); err != nil {
		return err
	}

	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		if _, err := h.AuthorizeModifyUser(ctx, pq, session, id); err != nil {
			return err
		}
		if err := pq.UpdateUser(ctx, id, data); err != nil {
			switch {
			case errors.Is(err, db.ErrNoRows):
				return services.NewNotFoundServiceError(err)
			}
			return err
		}
		return nil
	})
}

func (h *UserHandler) AuthorizeModifyUser(ctx context.Context, pq *db.PgxQueries, session *sessions.Session, userId *uuid.UUID) (*models.User, error) {
	sUserId := session.GetUserId()
	if sUserId == nil {
		return nil, services.NewUnauthenticatedServiceError(nil)
	}

	user, err := pq.GetUserForId(ctx, userId)
	if err != nil {
		// TODO: Handle database error
		return nil, err
	}

	if *sUserId == *userId || user.HasRole(models.USER_ROLE_ADMIN) {
		return user, nil
	}

	return nil, services.NewUnauthorizedServiceError(nil)
}
