package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (pq *PgxQueries) CreateApplication(ctx context.Context, businessId *uuid.UUID, postId int, userId *uuid.UUID) error {
	res, err := pq.tx.Exec(ctx, `
    INSERT INTO post_applications 
    (business_id, post_id, user_id) VALUES (@businessId, @postId, @userId)
    `, pgx.NamedArgs{
		"businessId": businessId,
		"postId":     postId,
		"userId":     userId,
	})

	if err != nil {
		return handlePgxError(err)
	}

	if res.RowsAffected() == 0 {
		return handlePgxError(ErrNoRows)
	}

	return nil
}
