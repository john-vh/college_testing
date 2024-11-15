package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/john-vh/college_testing/backend/models"
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

func (pq *PgxQueries) GetApplication(ctx context.Context, businessId *uuid.UUID, postId int, userId *uuid.UUID) (*models.UserApplication, error) {
	rows, err := pq.tx.Query(ctx, `
    SELECT post_applications.status, post_applications.created_at,
      json_build_object(
        'id', posts.id,
        'title', posts.title,
        'status', posts.status,
        'pay', posts.pay,
        'time_est', posts.time_est,
        'created_at', posts.created_at,
        'updated_at', posts.updated_at
    ) AS post,
      json_build_object(
        'id', businesses.id,
        'name', businesses.name,
        'status', businesses.status,
        'created_at', businesses.created_at
    ) AS business
    FROM post_applications
    LEFT JOIN posts ON posts.id = post_applications.post_id AND posts.business_id = post_applications.business_id
    LEFT JOIN businesses ON posts.business_id = businesses.id
    WHERE post_applications.post_id = @postId AND post_applications.business_id = @businessId AND post_applications.user_id = @userId
    `, pgx.NamedArgs{
		"userId":     userId,
		"businessId": businessId,
		"postId":     postId,
	})

	if err != nil {
		return nil, handlePgxError(err)
	}

	application, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[models.UserApplication])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return application, nil
}

func (pq *PgxQueries) SetApplicationStatus(ctx context.Context, businessId *uuid.UUID, postId int, userId *uuid.UUID, status models.ApplicationStatus) error {
	res, err := pq.tx.Exec(ctx, `
    UPDATE post_applications SET
    status = @status
    WHERE post_applications.post_id = @postId AND post_applications.business_id = @businessId AND post_applications.user_id = @userId
    `, pgx.NamedArgs{
		"businessId": businessId,
		"postId":     postId,
		"userId":     userId,
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
