package db

import (
	"context"
	"fmt"

	"github.com/WilliamTrojniak/StudentTests/backend/models"
	"github.com/WilliamTrojniak/StudentTests/backend/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (pq *PgxQueries) CreateUser(ctx context.Context, data *models.UserCreate) (*uuid.UUID, error) {
	userId, err := uuid.NewRandom()
	if err != nil {
		return nil, services.NewInternalServiceError(err)
	}

	_, err = pq.tx.Exec(ctx, `
    INSERT INTO users (id) VALUES (@userId)
    RETURNING id
    `,
		pgx.NamedArgs{
			"userId": userId,
		})

	if err != nil {
		return nil, handlePgxError(err)
	}

	return &userId, nil
}

func (pq *PgxQueries) SaveOpenIDAcct(ctx context.Context, openIDProvider string, data *models.OpenIDClaims) error {
	res, err := pq.tx.Exec(ctx, `
    INSERT INTO accounts 
    (provider, id, name, email, email_verified) VALUES (@provider, @id, @name, @email, @emailVerified)
    ON CONFLICT (provider, id) DO UPDATE
    SET (name, email, email_verified, updated_at) = (excluded.name, excluded.email, excluded.email_verified, NOW()) 
    `, pgx.NamedArgs{
		"provider":      openIDProvider,
		"id":            data.Id,
		"name":          data.Name,
		"email":         data.Email,
		"emailVerified": data.EmailVerified,
	})

	if err != nil {
		return handlePgxError(err)
	}

	if res.RowsAffected() == 0 {
		return handlePgxError(ErrNoRows)
	}

	return nil
}

func (pq *PgxQueries) LinkOpenIDAcct(ctx context.Context, openIDProvider string, acctData *models.OpenIDClaims, userId *uuid.UUID, isPrimary bool) error {

	res, err := pq.tx.Exec(ctx, `
    INSERT INTO user_accounts
    (user_id, account_provider, account_id, is_primary) VALUES (@userId, @accountProvider, @accountId, @isPrimary)
    `, pgx.NamedArgs{
		"userId":          userId,
		"accountProvider": openIDProvider,
		"accountId":       acctData.Id,
		"isPrimary":       isPrimary,
	})

	if err != nil {
		return handlePgxError(err)
	}

	if res.RowsAffected() == 0 {
		return handlePgxError(ErrNoRows)
	}

	return nil
}

func (pq *PgxQueries) GetLinkedUserId(ctx context.Context, provider string, accountId string) (*uuid.UUID, error) {
	row := pq.tx.QueryRow(ctx, `
    SELECT user_accounts.user_id
    FROM accounts
    LEFT JOIN user_accounts ON accounts.provider = user_accounts.account_provider AND accounts.id = user_accounts.account_id
    WHERE accounts.id = @accountId AND accounts.provider = @provider
    `,
		pgx.NamedArgs{
			"accountId": accountId,
			"provider":  provider,
		})

	var userId uuid.UUID
	err := row.Scan(&userId)
	if err != nil {
		return nil, handlePgxError(err)
	}
	if userId == uuid.Nil {
		return nil, nil
	}

	return &userId, nil

}

func (pq *PgxQueries) GetUserForId(ctx context.Context, id *uuid.UUID) (*models.User, error) {
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
    WHERE users.id = @userId
    `,
		pgx.NamedArgs{
			"userId": id,
		})

	/*

	   (SELECT array_remove(array_agg(account_info(accounts.provider, accounts.id, accounts.name, accounts.email, accounts.email_verified, user_accounts.is_primary, user_accounts.created_at, accounts.updated_at)), NULL)
	    FROM accounts
	    RIGHT JOIN user_accounts ON accounts.provider = user_accounts.account_provider AND accounts.id = user_accounts.account_id
	    WHERE user_accounts.user_id = @userId) AS accounts
	*/

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		fmt.Println(err)
		return nil, handlePgxError(err)
	}

	return user, nil
}
