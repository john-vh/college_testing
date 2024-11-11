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

func (pq *PgxQueries) GetBusinessOwner(ctx context.Context, businessId *uuid.UUID) (*models.User, error) {
	rows, err := pq.tx.Query(ctx, `
    SELECT 
      users.*, accounts.email, accounts.email_verified, accounts.name,
      (SELECT array_remove(array_agg(user_roles.role), NULL) 
       FROM user_roles WHERE user_roles.user_id = @userId) AS roles,
      (SELECT COALESCE(json_agg(accounts.*) FILTER (WHERE accounts.id IS NOT NULL), '[]')
       FROM user_accounts
       LEFT JOIN accounts ON user_accounts.account_provider = accounts.provider AND user_accounts.account_id = accounts.id 
        WHERE user_accounts.user_id = @userId
      ) as accounts
    FROM users
    LEFT JOIN user_accounts ON users.id = user_accounts.user_id AND user_accounts.is_primary = TRUE
    LEFT JOIN accounts ON user_accounts.account_provider = accounts.provider AND user_accounts.account_id = accounts.id
    LEFT JOIN businesses ON businesses.user_id = users.id
    WHERE businesses.id = @businessId
    `,
		pgx.NamedArgs{
			"businessId": businessId,
		})

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		fmt.Println(err)
		return nil, handlePgxError(err)
	}

	return user, nil
}

func (pq *PgxQueries) GetBusinesses(ctx context.Context, params *models.BusinessQueryParams) ([]models.Business, error) {
	if params == nil {
		params = &models.BusinessQueryParams{}
	}

	rows, err := pq.tx.Query(ctx, `
    SELECT * FROM businesses
    WHERE (@status::business_status IS NULL OR @status::business_status = businesses.status)
    AND (@userId::UUID IS NULL OR @userId::UUID = businesses.user_id)
    `,
		pgx.NamedArgs{
			"status": params.Status,
			"userId": params.UserId,
		})
	if err != nil {
		return nil, handlePgxError(err)
	}

	businesses, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Business])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return businesses, nil
}

func (pq *PgxQueries) GetBusinessForId(ctx context.Context, id *uuid.UUID) (*models.Business, error) {
	rows, err := pq.tx.Query(ctx, `
    SELECT * FROM businesses
    WHERE businesses.id = @businessId
    `,
		pgx.NamedArgs{
			"businessId": id,
		})
	if err != nil {
		return nil, handlePgxError(err)
	}

	business, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Business])
	if err != nil {
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
