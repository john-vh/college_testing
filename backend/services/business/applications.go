package business

import (
	"context"
	"errors"
	"fmt"

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

	return db.WithTx(ctx, h.store, func(pq *db.PgxQueries) error {
		sessionUser, err := pq.GetUserForId(ctx, sessionUserId)
		if err != nil {
			return services.NewUnauthorizedServiceError(err)
		}

		business, err := pq.GetBusinessForId(ctx, businessId)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return services.NewNotFoundServiceError(err)
			}
			return err
		}
		if business.Status != models.BUSINESS_STATUS_ACTIVE {
			return services.NewDataConflictServiceError(nil, "Business is not active")
		}

		post, err := pq.GetPostForId(ctx, businessId, postId)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return services.NewNotFoundServiceError(err)
			}
			return err
		}
		if post.Status != models.POST_STATUS_ACTIVE {
			return services.NewDataConflictServiceError(nil, "Post is not active")
		}
		targetUser, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return services.NewNotFoundServiceError(err)
			}
			return nil
		}

		if err := AuthorizeApplicationAction(sessionUser, APPLICATION_ACTION_CREATE, business, targetUser, nil, nil); err != nil {
			return err
		}

		if !targetUser.IsStudent() {
			if targetUser.Id == sessionUser.Id {
				return services.NewUnauthorizedServiceError(nil)
			}
			return services.NewDataConflictServiceError(err, "User is not a student")
		}

		err = pq.CreateApplication(ctx, businessId, postId, userId)
		if err != nil {
			if errors.Is(err, db.ErrUnique) {
				return services.NewDataConflictServiceError(err, "Application already exists")
			}
			return err
		}

		go func() {
			owner, err := db.WithTxRet(context.Background(), h.store, func(pq *db.PgxQueries) (*models.User, error) {
				return pq.GetBusinessOwner(context.Background(), businessId)
			})
			if err != nil {
				h.logger.Debug("Failed to get post owner while sending email")
				return
			}
			err = h.notifications.EnqueueWithTimeout(context.Background(), h.notifications.NewApplicationReceivedNotification(owner, targetUser, post))
			if err != nil {
				h.logger.Debug("Failed to send application confirmation email", "err", err)
			}
		}()
		return nil
	})
}

func (h *BusinessHandler) SetApplicationStatus(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, postId int, userId *uuid.UUID, status models.ApplicationStatus) error {
	h.logger.Debug("Setting application status", "business", businessId, "post", postId, "user", userId)
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
			if errors.Is(err, db.ErrNoRows) {
				return services.NewNotFoundServiceError(err)
			}
			return err
		}
		application, err := pq.GetApplication(ctx, businessId, postId, userId)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return services.NewNotFoundServiceError(err)
			}
			return err
		}
		targetUser, err := pq.GetUserForId(ctx, userId)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return services.NewNotFoundServiceError(err)
			}
			return err
		}

		var action ApplicationAction
		switch status {
		case models.APPLICATION_STATUS_ACCEPTED:
			action = APPLICATION_ACTION_ACCEPT
			if application.Status != models.APPLICATION_STATUS_PENDING {
				return services.NewDataConflictServiceError(nil, "Can not accept non-pending application")
			}
		case models.APPLICATION_STATUS_REJECTED:
			action = APPLICATION_ACTION_REJECT
			if application.Status != models.APPLICATION_STATUS_PENDING {
				return services.NewDataConflictServiceError(nil, "Can not reject non-pending application")
			}
		case models.APPLICATION_STATUS_COMPLETED:
			action = APPLICATION_ACTION_COMPLETE
			if !(application.Status == models.APPLICATION_STATUS_ACCEPTED || application.Status == models.APPLICATION_STATUS_INCOMPLETE) {
				return services.NewDataConflictServiceError(nil, "Can not complete non-accepted application")
			}
		case models.APPLICATION_STATUS_INCOMPLETE:
			action = APPLICATION_ACTION_INCOMPLETE
			if application.Status != models.APPLICATION_STATUS_ACCEPTED {
				return services.NewDataConflictServiceError(nil, "Can not mark non-accepted application incomplete")
			}
		case models.APPLICATION_STATUS_WITHDRAWN:
			action = APPLICATION_ACTION_WITHDRAW
			if !(application.Status == models.APPLICATION_STATUS_ACCEPTED || application.Status == models.APPLICATION_STATUS_PENDING) {
				return services.NewDataConflictServiceError(nil, "Can only withdraw pending and accepted applications.")
			}
		default:
			return services.NewBadRequestServiceError(fmt.Errorf("Invalid application status"))
		}

		if err := AuthorizeApplicationAction(sessionUser, action, business, targetUser, application, nil); err != nil {
			return err
		}

		err = pq.SetApplicationStatus(ctx, businessId, postId, userId, status)
		if err != nil {
			return err
		}
		// FIXME: Send to appropriate person based on application
		go db.WithTx(context.Background(), h.store, func(pq *db.PgxQueries) error {
			application, err := pq.GetApplication(context.Background(), businessId, postId, userId)
			if err != nil {
				h.logger.Debug("Failed to send application update email", "err", err)
				return err
			}
			err = h.notifications.EnqueueWithTimeout(context.Background(), h.notifications.NewApplicationUpdatedNotification(targetUser, targetUser, application))
			if err != nil {
				h.logger.Debug("Failed to send application update email", "err", err)
				return err
			}
			return nil
		})
		return nil
	})
	return err
}

func (h *BusinessHandler) GetPostApplications(ctx context.Context, session *sessions.Session, businessId *uuid.UUID, postId int) (*models.PostApplications, error) {
	h.logger.Debug("Retrieving post applications", "Business Id", businessId, "Post Id", postId)
	userId := session.GetUserId()
	if userId == nil {
		return nil, services.NewUnauthenticatedServiceError(nil)
	}
	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) (*models.PostApplications, error) {
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

		applications, err := pq.GetApplicationsForPost(ctx, businessId, postId)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return nil, services.NewNotFoundServiceError(err)
			}
			return nil, err
		}
		return applications, nil
	})
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

	return db.WithTxRet(ctx, h.store, func(pq *db.PgxQueries) ([]models.UserApplication, error) {
		user, err := pq.GetUserForId(ctx, params.UserId)
		if err != nil {
			return nil, services.NewUnauthorizedServiceError(err)
		}
		if err := AuthorizeApplicationAction(user, APPLICATION_ACTION_READ_USER, nil, nil, nil, params); err != nil {
			return nil, err
		}

		applications, err := pq.GetUserApplications(ctx, params)
		if err != nil {
			return nil, err
		}
		return applications, nil
	})
}

type ApplicationAction string

const (
	APPLICATION_ACTION_CREATE     ApplicationAction = "application:create"
	APPLICATION_ACTION_READ_USER  ApplicationAction = "application:read_user"
	APPLICATION_ACTION_READ       ApplicationAction = "application:read"
	APPLICATION_ACTION_REJECT     ApplicationAction = "application:reject"
	APPLICATION_ACTION_ACCEPT     ApplicationAction = "application:accept"
	APPLICATION_ACTION_COMPLETE   ApplicationAction = "application:complete"
	APPLICATION_ACTION_INCOMPLETE ApplicationAction = "application:incomplete"
	APPLICATION_ACTION_WITHDRAW   ApplicationAction = "application:withdraw"
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
			case APPLICATION_ACTION_ACCEPT:
				return nil
			case APPLICATION_ACTION_REJECT:
				return nil
			case APPLICATION_ACTION_COMPLETE:
				return nil
			case APPLICATION_ACTION_INCOMPLETE:
				return nil
			case APPLICATION_ACTION_WITHDRAW:
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
			case APPLICATION_ACTION_ACCEPT:
				if business != nil && business.UserId == user.Id {
					return nil
				}
			case APPLICATION_ACTION_REJECT:
				if business != nil && business.UserId == user.Id {
					return nil
				}
			case APPLICATION_ACTION_COMPLETE:
				if business != nil && business.UserId == user.Id {
					return nil
				}
			case APPLICATION_ACTION_INCOMPLETE:
				if business != nil && business.UserId == user.Id {
					return nil
				}
			case APPLICATION_ACTION_WITHDRAW:
				if targetUser != nil && targetUser.Id == user.Id {
					return nil
				}
			}
		}
	}

	return services.NewUnauthorizedServiceError(nil)
}
