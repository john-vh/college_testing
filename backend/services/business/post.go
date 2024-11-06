package business

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services/notifications"
	"github.com/john-vh/college_testing/backend/services/sessions"
)

func (h *BusinessHandler) GetPosts(ctx context.Context, session *sessions.Session, params *models.PostQueryParams) ([]models.Post, error) {
	h.logger.Debug("Retreiving posts")
	// TODO: Authorize session to retreive posts

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.Post, error) {
		return pq.GetPosts(ctx, params)
	})
}

func (h *BusinessHandler) CreatePost(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, data *models.PostCreate) (*models.Post, error) {
	h.logger.Debug("Creating post", "Business Id", businessId)
	// TODO: Authorize session to modify business

	if err := models.ValidateData(data); err != nil {
		return nil, err
	}

	post, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.Post, error) {
		return pq.CreatePost(ctx, businessId, data)
	})

	if err != nil {
		h.logger.Debug("Error creating post", "err", err)
		return nil, err
	}
	h.logger.Debug("Created post", "id", post.Id)
	return post, nil
}

func (h *BusinessHandler) UpdatePost(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, postId int, data *models.PostUpdate) error {
	h.logger.Debug("Updating post", "Business Id", businessId, "Post Id", postId)
	// TODO: Authorize session to modify post

	if err := models.ValidateData(data); err != nil {
		return err
	}

	err := db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		return pq.UpdatePost(ctx, businessId, postId, data)
	})
	if err != nil {
		h.logger.Debug("Error updating post", "err", err)
		return err
	}
	h.logger.Debug("Updated post", "Business Id", businessId, "Post Id", postId)
	return nil
}

func (h *BusinessHandler) SetPostStatus(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, postId int, status models.PostStatus) error {
	h.logger.Debug("Setting post status", "Business Id", businessId, "Post Id", postId, "status", status.String())
	// TODO: Authorize session to modify post

	err := db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		return pq.SetPostStatus(ctx, businessId, postId, status)
	})
	if err != nil {
		h.logger.Debug("Error setting post status", "err", err)
		return err
	}
	h.logger.Debug("Set post status", "Business Id", businessId, "Post Id", postId)
	return nil
}

func (h *BusinessHandler) CreateApplication(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, postId int, userId *uuid.UUID) error {
	h.logger.Debug("Creating application", "Business Id", businessId, "Post Id", postId, "User Id", userId)
	// TODO: Authorize session to apply to post

	err := db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		return pq.CreateApplication(ctx, businessId, postId, userId)
	})
	if err != nil {
		h.logger.Debug("Error creating application", "err", err)
		return err
	}

	go func() {
		user, err := db.WithTxRet(context.TODO(), h.store, func(pq *db.PgxQueries) (*models.User, error) {
			user, err := pq.GetUserForId(context.TODO(), userId)
			return user, err
		})
		if err != nil {
			h.logger.Debug("failed to get user while sending email")
			return
		}

		email := &notifications.MailInfo{
			ToList:  []string{user.Email},
			Subject: "Application Received",
			Body: fmt.Sprintf(
				`Hi %v!\n
        This is a confirmation email for applying! Expect to hear back from them soon.
        `, user.Name),
		}
		err = h.notifications.SendMsg(email.ToList, email)
		if err != nil {
			h.logger.Debug("Failed to send application confirmation email")
		}
	}()

	h.logger.Debug("Created application", "Business Id", businessId, "Post Id", postId, "User Id", userId)
	return nil
}

func (h *BusinessHandler) GetPostApplications(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, postId int) (*models.PostApplications, error) {
	h.logger.Debug("Retrieving applications", "Business Id", businessId, "Post Id", postId)
	// TODO: Authorize session to get applications

	applications, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.PostApplications, error) {
		return pq.GetApplicationsForPost(ctx, businessId, postId)
	})
	if err != nil {
		h.logger.Debug("Error retrieving applications", "err", err)
		return nil, err
	}

	h.logger.Debug("Retrieved applications", "Business Id", businessId, "Post Id", postId)
	return applications, nil
}
