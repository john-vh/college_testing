package db

import (
	"context"

	"github.com/WilliamTrojniak/StudentTests/backend/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (pq *PgxQueries) CreatePost(ctx context.Context, businessId *uuid.UUID, data *models.PostCreate) (*models.Post, error) {
	rows, err := pq.tx.Query(ctx, `
    INSERT INTO posts 
    (business_id, title, description) VALUES (@businessId, @title, @description)
    RETURNING posts.*
    `, pgx.NamedArgs{
		"businessId":  businessId,
		"title":       data.Title,
		"description": data.Desc,
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
    (title, description, updated_at) = (@title, @description, NOW())
    WHERE posts.id = @postId AND posts.business_id = @businessId
    `, pgx.NamedArgs{
		"businessId":  businessId,
		"postId":      postId,
		"title":       data.Title,
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
