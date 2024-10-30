package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services"
)

func (pq *PgxQueries) CreateBusiness(ctx context.Context, ownerId *uuid.UUID, data *models.BusinessCreate) (*models.Business, error) {
	businessId, err := uuid.NewRandom()
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}

	rows, err := pq.tx.Query(ctx, `
    INSERT INTO businesses
    (user_id, id, name, website, description) VALUES (@userId, @businessId, @name, @website, @description)
    RETURNING businesses.*
    `, pgx.NamedArgs{
		"userId":      ownerId,
		"businessId":  businessId,
		"name":        data.Name,
		"website":     data.Website,
		"description": data.Desc,
	})

	if err != nil {
		return nil, handlePgxError(err)
	}

	business, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Business])
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, handlePgxError(err)
	}

	return business, nil
}

func (pq *PgxQueries) UpdateBusiness(ctx context.Context, businessId *uuid.UUID, data *models.BusinessUpdate) error {

	res, err := pq.tx.Exec(ctx, `
    UPDATE businesses SET
    (website, description) = (@website, @description)
    WHERE businesses.id = @businessId
    `, pgx.NamedArgs{
		"businessId":  businessId,
		"website":     data.Website,
		"description": data.Desc,
	})

	if err != nil {
		return handlePgxError(err)
	}

	if res.RowsAffected() == 0 {
		return handlePgxError(ErrNoRows)
	}

	return nil
}

func (pq *PgxQueries) SetBusinessStatus(ctx context.Context, businessId *uuid.UUID, status models.BusinessStatus) error {
	res, err := pq.tx.Exec(ctx, `
    UPDATE businesses SET
    status = @status
    WHERE businesses.id = @businessId
    `, pgx.NamedArgs{
		"businessId": businessId,
		"status":     status,
	})

	if err != nil {
		return handlePgxError(err)
	}

	if res.RowsAffected() == 0 {
		return handlePgxError(ErrNoRows)
	}

	return nil
}
