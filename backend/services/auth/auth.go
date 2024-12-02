package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/models"
	"github.com/john-vh/college_testing/backend/services"
	"github.com/john-vh/college_testing/backend/services/sessions"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	logger          *slog.Logger
	sessions        *sessions.SessionsHandler
	store           *db.PgxStore
	providers       map[string]providerConfig
	handleErr       services.ServicesHTTPErrorHandler
	redirectBaseURI string
}

type providerConfig struct {
	config   oauth2.Config
	provider *oidc.Provider
}

type ProviderConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	Scopes       []string
}

func NewAuthHandler(logger *slog.Logger, errHandler services.ServicesHTTPErrorHandler, sessions *sessions.SessionsHandler, store *db.PgxStore, baseURL string, redirectBaseURL string, providerConfigs map[string]ProviderConfig) (*AuthHandler, error) {
	providers := make(map[string]providerConfig)
	for provider, config := range providerConfigs {
		providerClient, err := oidc.NewProvider(context.TODO(), config.Issuer)
		if err != nil {
			return nil, err
		}
		config := oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Endpoint:     providerClient.Endpoint(),
			RedirectURL:  fmt.Sprintf("%v/auth/%v/callback", baseURL, provider),
			Scopes:       config.Scopes,
		}
		providers[provider] = providerConfig{provider: providerClient, config: config}
	}

	return &AuthHandler{
		logger:          logger,
		handleErr:       errHandler,
		sessions:        sessions,
		store:           store,
		providers:       providers,
		redirectBaseURI: redirectBaseURL,
	}, nil
}

func (auth *AuthHandler) SaveAccount(ctx context.Context, openidProvider string, acctData *models.OpenIDClaims, userData *models.UserCreate) (*uuid.UUID, error) {
	auth.logger.Info("Auth saving account", "createUser", userData != nil)
	return db.WithTxRet(ctx, auth.store, func(pq *db.PgxQueries) (*uuid.UUID, error) {
		err := pq.SaveOpenIDAcct(ctx, openidProvider, acctData)
		if err != nil {
			auth.logger.Debug("Failed to save openID account")
			return nil, err
		}

		linkedUserId, err := pq.GetLinkedUserId(ctx, openidProvider, acctData.Id)
		if err != nil {
			auth.logger.Debug("Failed to get linked user id")
			return nil, err
		}

		// Return the already associated user, or do not create a user if no data is given
		if linkedUserId != nil || userData == nil {
			auth.logger.Debug("Skipping user create", "linked_id", linkedUserId)
			return linkedUserId, nil
		}

		userId, err := pq.CreateUser(ctx, userData)
		if err != nil {
			auth.logger.Debug("Failed to create user")
			return nil, err
		}
		err = pq.LinkOpenIDAcct(ctx, openidProvider, acctData, userId, true)
		if err != nil {
			auth.logger.Debug("Failed to link account to user")
			return nil, err
		}

		return userId, nil
	})
}

func (auth *AuthHandler) GetLinkedUser(ctx context.Context, openidProvider string, acctData *models.OpenIDClaims) (*uuid.UUID, error) {
	return db.WithTxRet(ctx, auth.store, func(pq *db.PgxQueries) (*uuid.UUID, error) {
		linkedUserId, err := pq.GetLinkedUserId(ctx, openidProvider, acctData.Id)
		if err != nil {
			return nil, err
		}
		return linkedUserId, nil
	})
}

func (auth *AuthHandler) LinkAccount(ctx context.Context, userId *uuid.UUID, openidProvider string, acctData *models.OpenIDClaims) error {
	auth.logger.Info("Auth linking account", "account_provider", openidProvider, "account_id", acctData.Id, "user_id", userId)
	return db.WithTx(ctx, auth.store, func(pq *db.PgxQueries) error {
		err := pq.SaveOpenIDAcct(ctx, openidProvider, acctData)
		if err != nil {
			return err
		}

		linkedUserId, err := pq.GetLinkedUserId(ctx, openidProvider, acctData.Id)
		if err != nil {
			return err
		}

		if linkedUserId != nil && (*linkedUserId) == (*userId) {
			return nil
		} else if linkedUserId != nil {
			return services.NewDataConflictServiceError(nil, "Account is already linked to another user")
		}

		err = pq.LinkOpenIDAcct(ctx, openidProvider, acctData, userId, false)
		if err != nil {
			return err
		}
		return nil
	})
}

func setCallbackCookie(w http.ResponseWriter, name string, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
