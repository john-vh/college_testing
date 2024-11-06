package business

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services"
	"github.com/john-vh/college_testing/backend/services/sessions"
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

func (h *BusinessHandler) GetBusinesses(ctx context.Context, session *sessions.Session, params *models.BusinessQueryParams) ([]models.Business, error) {
	// TODO: Authorization of session to get the requested businesses
	// Need to consider the params

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.Business, error) {
		return pq.GetBusinesses(ctx, params)
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
