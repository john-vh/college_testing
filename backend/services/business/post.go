package business

import (
	"context"

	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services"
	"github.com/john-vh/college_testing/backend/services/sessions"
)

func (h *BusinessHandler) GetPosts(ctx context.Context, session *sessions.Session, params *models.PostQueryParams) ([]models.Post, error) {
	h.logger.Debug("Retreiving posts")
	userId := session.GetUserId()
	if userId == nil {
		return nil, services.NewUnauthorizedServiceError(nil)
	}

	if params == nil {
		params = &models.PostQueryParams{}
	}

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.Post, error) {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return nil, services.NewUnauthorizedServiceError(err)
		}

		var business *models.Business
		if params.BusinessId != nil {
			business, err = pq.GetBusinessForId(ctx, params.BusinessId)
			if err != nil {
				// HACK: Need to implement better error handling between service and db levels
				business = nil
			}
		}
		if err := AuthorizePostAction(user, POST_ACTION_READ, business, nil, params); err != nil {
			return nil, err
		}

		return pq.GetPosts(ctx, params)
	})
}

func (h *BusinessHandler) CreatePost(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, data *models.PostCreate) (*models.Post, error) {
	h.logger.Debug("Creating post", "Business Id", businessId)
	userId := session.GetUserId()
	if userId == nil {
		return nil, services.NewUnauthenticatedServiceError(nil)
	}

	if err := models.ValidateData(data); err != nil {
		return nil, err
	}
	post, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.Post, error) {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return nil, services.NewUnauthorizedServiceError(err)
		}
		business, err := pq.GetBusinessForId(ctx, businessId)
		if err != nil {
			return nil, services.NewUnauthorizedServiceError(err)
		}

		if err := AuthorizePostAction(user, POST_ACTION_CREATE, business, nil, nil); err != nil {
			return nil, err
		}

		if business.Status != models.BUSINESS_STATUS_ACTIVE {
			return nil, services.NewDataConflictServiceError(nil, "Business is not active")
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
	userId := session.GetUserId()
	if userId == nil {
		return services.NewUnauthenticatedServiceError(nil)
	}

	h.logger.Debug("Updating post", "Business Id", businessId, "Post Id", postId)
	err := db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return services.NewUnauthorizedServiceError(err)
		}
		business, err := pq.GetBusinessForId(ctx, businessId)
		if err != nil {
			return err
		}
		post, err := pq.GetPostForId(ctx, businessId, postId)
		if err != nil {
			return err
		}

		if err := AuthorizePostAction(user, POST_ACTION_UPDATE, business, post, nil); err != nil {
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
	userId := session.GetUserId()
	if userId == nil {
		return services.NewUnauthenticatedServiceError(nil)
	}

	err := db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return services.NewUnauthorizedServiceError(err)
		}
		business, err := pq.GetBusinessForId(ctx, businessId)
		if err != nil {
			return err
		}
		post, err := pq.GetPostForId(ctx, businessId, postId)
		if err != nil {
			return err
		}

		if err := AuthorizePostAction(user, POST_ACTION_UPDATE, business, post, nil); err != nil {
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

type PostAction string

const (
	POST_ACTION_CREATE PostAction = "post:create"
	POST_ACTION_UPDATE PostAction = "post:update"
	POST_ACTION_READ   PostAction = "post:read"
)

func AuthorizePostAction(user *models.User, action PostAction, business *models.Business, post *models.Post, query *models.PostQueryParams) error {
	if user == nil {
		return services.NewUnauthenticatedServiceError(nil)
	}

	for _, role := range user.Roles {
		switch role {
		case models.USER_ROLE_ADMIN:
			switch action {
			case POST_ACTION_CREATE:
				return nil
			case POST_ACTION_UPDATE:
				return nil
			case POST_ACTION_READ:
				return nil
			}
		case models.USER_ROLE_USER:
			switch action {
			case POST_ACTION_CREATE:
				if business != nil && (business.UserId == user.Id) {
					return nil
				}
			case POST_ACTION_UPDATE:
				if business != nil && post != nil &&
					(business.UserId == user.Id && business.Id == post.BusinessId) {
					return nil
				}
			case POST_ACTION_READ:
				if query != nil && ((query.UserId != nil && *query.UserId == user.Id) ||
					(query.Status != nil && *query.Status == models.POST_STATUS_ACTIVE) ||
					(business != nil && query.BusinessId != nil && business.Id == *query.BusinessId && business.UserId == user.Id)) {
					return nil
				}
			}
		}
	}

	return services.NewUnauthorizedServiceError(nil)
}
