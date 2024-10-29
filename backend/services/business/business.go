package business

import (
	"context"
	"log/slog"

	"github.com/WilliamTrojniak/StudentTests/backend/db"
	"github.com/WilliamTrojniak/StudentTests/backend/models"
	"github.com/WilliamTrojniak/StudentTests/backend/services"
	"github.com/WilliamTrojniak/StudentTests/backend/services/sessions"
	"github.com/google/uuid"
)

type BusinessHandler struct {
	logger    *slog.Logger
	sessions  *sessions.SessionsHandler
	store     *db.PgxStore
	handleErr services.ServicesHTTPErrorHandler
}

func NewBusinessHandler(
	logger *slog.Logger,
	errHandler services.ServicesHTTPErrorHandler,
	sessions *sessions.SessionsHandler,
	store *db.PgxStore,
) *BusinessHandler {
	return &BusinessHandler{
		logger:    logger,
		sessions:  sessions,
		store:     store,
		handleErr: errHandler,
	}
}

func (h *BusinessHandler) RequestBusiness(ctx context.Context, session *sessions.Session, ownerId *uuid.UUID, data *models.BusinessCreate) (*models.Business, error) {
	// TODO: Authorization of the session to modify the requested user

	err := models.ValidateData(data)
	if err != nil {
		return nil, err
	}

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.Business, error) {
		return pq.CreateBusiness(ctx, ownerId, data)
	})
}

func (h *BusinessHandler) UpdateBusiness(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, data *models.BusinessUpdate) error {
	// TODO: Authorization of the session to modify the requested business

	err := models.ValidateData(data)
	if err != nil {
		return err
	}

	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		return pq.UpdateBusiness(ctx, businessId, data)
	})

}

func (h *BusinessHandler) ApproveBusiness(ctx context.Context, session *sessions.Session, businessId *uuid.UUID) error {
	// TODO: Authorization of the session to modify the requested business

	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		return pq.SetBusinessStatus(ctx, businessId, models.BUSINESS_STATUS_ACTIVE)
	})
}
