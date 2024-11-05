package user

import (
	"context"
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
	err := h.AuthorizeModifyUser(ctx, session, id)
	if err != nil {
		return nil, err
	}

	user, err := h.getUserById(ctx, id)
	if err != nil {
		h.logger.Debug("Error getting user", "err", err)
		return nil, err
	}

	return user, nil
}

func (h *UserHandler) getUserById(ctx context.Context, id *uuid.UUID) (*models.User, error) {
	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.User, error) {
		return pq.GetUserForId(ctx, id)
	})
}

func (h *UserHandler) AuthorizeModifyUser(ctx context.Context, session *sessions.Session, userId *uuid.UUID) (err error) {
	sUserId := session.GetUserId()
	if sUserId == nil {
		return services.NewUnauthorizedServiceError(nil)
	}

	if *sUserId == *userId {
		return nil
	}

	user, err := h.getUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user.HasRole(models.USER_ROLE_ADMIN) {
		return nil
	}

	return services.NewUnauthorizedServiceError(nil)
}
