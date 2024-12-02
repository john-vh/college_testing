package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/john-vh/college_testing/backend/models"
)

func (pq *PgxQueries) GetPosts(ctx context.Context, params *models.PostQueryParams) ([]models.Post, error) {
	rows, err := pq.tx.Query(ctx, `
    SELECT posts.*
    FROM posts
    LEFT JOIN businesses ON businesses.id = posts.business_id
    LEFT JOIN users ON businesses.user_id = users.id
    WHERE (@status::post_status IS NULL OR @status::post_status = posts.status)
    AND (@businessId::UUID IS NULL OR @businessId::UUID = posts.business_id)
    AND (@userId::UUID IS NULL OR @userId::UUID = users.id)
    AND businesses.status = @businessActive
    `, pgx.NamedArgs{
		"status":         params.Status,
		"businessId":     params.BusinessId,
		"userId":         params.UserId,
		"businessActive": models.BUSINESS_STATUS_ACTIVE,
	})
	if err != nil {
		return nil, handlePgxError(err)
	}

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Post])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return posts, nil
}

func (pq *PgxQueries) GetPostForId(ctx context.Context, businessId *uuid.UUID, postId int) (*models.Post, error) {
	rows, err := pq.tx.Query(ctx, `
    SELECT posts.*
    FROM posts
    WHERE posts.business_id = @businessId AND posts.id = @postId
    `, pgx.NamedArgs{
		"businessId": businessId,
		"postId":     postId,
	})
	if err != nil {
		return nil, handlePgxError(err)
	}

	post, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[models.Post])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return post, nil
}

func (pq *PgxQueries) CreatePost(ctx context.Context, businessId *uuid.UUID, data *models.PostCreate) (*models.Post, error) {
	rows, err := pq.tx.Query(ctx, `
    INSERT INTO posts 
    (business_id, title, description, pay, time_est) 
    VALUES (@businessId, @title, @description, @pay, @timeEst)
    RETURNING posts.*
    `, pgx.NamedArgs{
		"businessId":  businessId,
		"title":       data.Title,
		"description": data.Desc,
		"pay":         data.Pay,
		"timeEst":     data.TimeEst,
	})

	if err != nil {
		return nil, handlePgxError(err)
	}

	post, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Post])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return post, nil
}

func (pq *PgxQueries) UpdatePost(ctx context.Context, businessId *uuid.UUID, postId int, data *models.PostUpdate) error {
	res, err := pq.tx.Exec(ctx, `
    UPDATE posts SET
    (title, description, pay, time_est, updated_at) = (@title, @description, @pay, @time_est, NOW())
    WHERE posts.id = @postId AND posts.business_id = @businessId
    `, pgx.NamedArgs{
		"businessId":  businessId,
		"postId":      postId,
		"title":       data.Title,
		"description": data.Desc,
		"pay":         data.Pay,
		"timeEst":     data.TimeEst,
	})

	if err != nil {
		return handlePgxError(err)
	}

	if res.RowsAffected() == 0 {
		return handlePgxError(ErrNoRows)
	}

	return nil
}

func (pq *PgxQueries) SetPostStatus(ctx context.Context, businessId *uuid.UUID, postId int, status models.PostStatus) error {
	// Does not affected updated-at time
	res, err := pq.tx.Exec(ctx, `
    UPDATE posts SET
    status = @status
    WHERE posts.id = @postId AND posts.business_id = @businessId
    `, pgx.NamedArgs{
		"businessId": businessId,
		"postId":     postId,
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

func (pq *PgxQueries) GetApplicationsForPost(ctx context.Context, businessId *uuid.UUID, postId int) (*models.PostApplications, error) {
	rows, err := pq.tx.Query(ctx, `
    SELECT post_applications.notes, post_applications.status, post_applications.created_at,
      json_build_object(
      'id', users.id,
      'created_at', users.created_at,
      'email', accounts.email,
      'name', accounts.name,
      'email_verified', accounts.email_verified,
      'status', users.status
    ) AS user
    FROM post_applications
    LEFT JOIN users on post_applications.user_id = users.id 
    LEFT JOIN user_accounts ON users.id = user_accounts.user_id
    LEFT JOIN accounts ON user_accounts.account_provider = accounts.provider AND user_accounts.account_id = accounts.id
    WHERE post_applications.business_id = @businessId AND post_applications.post_id = @postId AND user_accounts.is_primary = TRUE
    `, pgx.NamedArgs{
		"businessId": businessId,
		"postId":     postId,
	})

	if err != nil {
		return nil, handlePgxError(err)
	}

	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.PostApplicationData])
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, handlePgxError(err)
	}

	applications := &models.PostApplications{
		BusinessId:   *businessId,
		PostId:       postId,
		Applications: data,
	}
	return applications, nil
}

func (pq *PgxQueries) GetUserApplications(ctx context.Context, params *models.UserApplicationQueryParams) ([]models.UserApplication, error) {
	if params == nil {
		params = &models.UserApplicationQueryParams{}
	}

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
    WHERE (@userId::UUID IS NULL OR @userId = post_applications.user_id)
    AND (@applicationStatus::post_application_status IS NULL OR @applicationStatus = post_applications.status)
    AND (@postStatus::post_status IS NULL OR @postStatus = posts.status)
    `, pgx.NamedArgs{
		"userId":            params.UserId,
		"applicationStatus": params.ApplicationStatus,
		"postStatus":        params.PostStatus,
	})

	if err != nil {
		return nil, handlePgxError(err)
	}

	applications, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.UserApplication])
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, handlePgxError(err)
	}

	return applications, nil
}
