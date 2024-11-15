package business

import (
	"context"

	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services"
	"github.com/john-vh/college_testing/backend/services/sessions"
)

func (h *BusinessHandler) CreateApplication(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, postId int, userId *uuid.UUID) error {
	h.logger.Debug("Creating application", "Business Id", businessId, "Post Id", postId, "User Id", userId)
	sessionUserId := session.GetUserId()
	if sessionUserId == nil {
		return services.NewUnauthenticatedServiceError(nil)
	}

	err := db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		sessionUser, err := pq.GetUserForId(ctx, sessionUserId)
		if err != nil {
			return services.NewUnauthorizedServiceError(err)
		}

		business, err := pq.GetBusinessForId(ctx, businessId)
		if err != nil {
			return err
		}
		if business.Status != models.BUSINESS_STATUS_ACTIVE {
			return services.NewDataConflictServiceError(nil, "Business is not active")
		}

		post, err := pq.GetPostForId(ctx, businessId, postId)
		if err != nil {
			return err
		}
		if post.Status != models.POST_STATUS_ACTIVE {
			return services.NewDataConflictServiceError(nil, "Post is not active")
		}
		targetUser, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return nil
		}

		if err := AuthorizeApplicationAction(sessionUser, APPLICATION_ACTION_CREATE, business, targetUser, nil, nil); err != nil {
			return err
		}

		if !targetUser.IsStudent() {
			return services.NewDataConflictServiceError(err, "User is not a student")
		}

		return pq.CreateApplication(ctx, businessId, postId, userId)
	})
	if err != nil {
		h.logger.Debug("Error creating application", "err", err)
		return err
	}

	h.logger.Debug("Created application", "Business Id", businessId, "Post Id", postId, "User Id", userId)
	return nil
}

func (h *BusinessHandler) GetPostApplications(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, postId int) (*models.PostApplications, error) {
	h.logger.Debug("Retrieving post applications", "Business Id", businessId, "Post Id", postId)
	userId := session.GetUserId()
	if userId == nil {
		return nil, services.NewUnauthenticatedServiceError(nil)
	}
	applications, err := db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.PostApplications, error) {
		user, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			return nil, services.NewUnauthorizedServiceError(err)
		}
		business, err := pq.GetBusinessForId(ctx, businessId)
		if err != nil {
			return nil, services.NewUnauthorizedServiceError(err)
		}

		if err := AuthorizeApplicationAction(user, APPLICATION_ACTION_READ, business, nil, nil, nil); err != nil {
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
		user, err := pq.GetUserForId(ctx, params.UserId)
		if err != nil {
			return nil, services.NewUnauthorizedServiceError(err)
		}
		if err := AuthorizeApplicationAction(user, APPLICATION_ACTION_READ_USER, nil, nil, nil, params); err != nil {
			return nil, err
		}

		return pq.GetUserApplications(ctx, params)
	})

	if err != nil {
		return nil, err
	}

	return applications, nil
}

type ApplicationAction string

const (
	APPLICATION_ACTION_CREATE    ApplicationAction = "application:create"
	APPLICATION_ACTION_READ_USER ApplicationAction = "application:read_user"
	APPLICATION_ACTION_READ      ApplicationAction = "application:read"
)

func AuthorizeApplicationAction(user *models.User, action ApplicationAction, business *models.Business, targetUser *models.User, application *models.UserApplication, query *models.UserApplicationQueryParams) error {
	if user == nil {
		return services.NewUnauthenticatedServiceError(nil)
	}

	for _, role := range user.Roles {
		switch role {
		case models.USER_ROLE_ADMIN:
			switch action {
			case APPLICATION_ACTION_CREATE:
				return nil
			case APPLICATION_ACTION_READ_USER:
				return nil
			case APPLICATION_ACTION_READ:
				return nil
			}
		case models.USER_ROLE_USER:
			switch action {
			case APPLICATION_ACTION_CREATE:
				if targetUser != nil && (targetUser.Id == user.Id) {
					return nil
				}
			case APPLICATION_ACTION_READ_USER:
				if query != nil && (*query.UserId == user.Id) {
					return nil
				}
			case APPLICATION_ACTION_READ:
				if business != nil && business.UserId == user.Id {
					return nil
				}
			}
		}
	}

	return services.NewUnauthorizedServiceError(nil)
}
