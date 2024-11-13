package business

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services"
	"github.com/john-vh/college_testing/backend/services/notifications"
	"github.com/john-vh/college_testing/backend/services/sessions"
)

func (h *BusinessHandler) GetPosts(ctx context.Context, session *sessions.Session, params *models.PostQueryParams) ([]models.Post, error) {
	h.logger.Debug("Retreiving posts")
	userId := session.GetUserId()
	if userId == nil {
		return nil, services.NewUnauthorizedServiceError(nil)
	}

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.Post, error) {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return nil, services.NewUnauthorizedServiceError(err)
		}

		if params == nil {
			params = &models.PostQueryParams{}
		}

		var business *models.Business
		if params.BusinessId != nil {
			business, err = pq.GetBusinessForId(ctx, params.BusinessId)
			if err != nil {
				// HACK: Need to implement better error handling between service and db levels
				business = nil
			}
		}

		if (params.Status == nil || *params.Status != models.POST_STATUS_ACTIVE) &&
			!(user.HasRole(models.USER_ROLE_ADMIN) ||
				(params.UserId != nil && *params.UserId == *userId) ||
				(business != nil && business.UserId == *userId)) {
			return nil, services.NewUnauthorizedServiceError(fmt.Errorf("Attempted to view inactive posts"))
		}

		return pq.GetPosts(ctx, params)
	})
}

func (h *BusinessHandler) CreatePost(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, data *models.PostCreate) (*models.Post, error) {
	h.logger.Debug("Creating post", "Business Id", businessId)

	if err := models.ValidateData(data); err != nil {
		return nil, err
	}

	post, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.Post, error) {
		if _, err := h.AuthorizeModifyBusiness(ctx, pq, session, businessId); err != nil {
			return nil, err
		}
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
	err := db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		if _, err := h.AuthorizeModifyBusiness(ctx, pq, session, businessId); err != nil {
			return err
		}
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
	h.logger.Debug("Setting post status", "Business Id", businessId, "Post Id", postId, "status", status)
	err := db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		if _, err := h.AuthorizeModifyBusiness(ctx, pq, session, businessId); err != nil {
			return err
		}
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

	err := db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		user, err := h.users.AuthorizeModifyUser(ctx, pq, session, userId)
		if err != nil {
			return err
		}
		if !user.IsStudent() {
			return services.NewDataConflictServiceError(err, "User is not a student")
		}
		post, err := pq.GetPostForId(ctx, businessId, postId)
		if err != nil {
			return err
		}
		if post.Status != models.POST_STATUS_ACTIVE {
			return services.NewNotFoundServiceError(nil)
		}

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
	h.logger.Debug("Retrieving post applications", "Business Id", businessId, "Post Id", postId)
	applications, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.PostApplications, error) {
		if _, err := h.AuthorizeModifyBusiness(ctx, pq, session, businessId); err != nil {
			return nil, err
		}
		return pq.GetApplicationsForPost(ctx, businessId, postId)
	})
	if err != nil {
		h.logger.Debug("Error retrieving applications", "err", err)
		return nil, err
	}

	h.logger.Debug("Retrieved applications", "Business Id", businessId, "Post Id", postId)
	return applications, nil
}

func (h *BusinessHandler) GetUserApplications(ctx context.Context, session *sessions.Session, params *models.UserApplicationQueryParams) ([]models.UserApplication, error) {
	h.logger.Debug("Retreiving user applications.")
	userId := session.GetUserId()
	if userId == nil {
		return nil, services.NewUnauthenticatedServiceError(nil)
	}

	if params == nil {
		params = &models.UserApplicationQueryParams{}
	}

	applications, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.UserApplication, error) {
		if params.UserId == nil || *params.UserId != *userId {
			user, err := pq.GetUserForId(ctx, params.UserId)
			if err != nil {
				return nil, err
			}
			if !user.HasRole(models.USER_ROLE_ADMIN) {
				return nil, services.NewUnauthorizedServiceError(nil)
			}
		}

		return pq.GetUserApplications(ctx, params)
	})

	if err != nil {
		return nil, err
	}

	return applications, nil
}
