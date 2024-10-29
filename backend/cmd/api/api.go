package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/WilliamTrojniak/StudentTests/backend/cache"
	"github.com/WilliamTrojniak/StudentTests/backend/db"
	"github.com/WilliamTrojniak/StudentTests/backend/env"
	"github.com/WilliamTrojniak/StudentTests/backend/services"
	"github.com/WilliamTrojniak/StudentTests/backend/services/auth"
	"github.com/WilliamTrojniak/StudentTests/backend/services/business"
	"github.com/WilliamTrojniak/StudentTests/backend/services/sessions"
	"github.com/WilliamTrojniak/StudentTests/backend/services/user"
	"github.com/redis/go-redis/v9"
)

type APIServer struct {
	addr  string
	cfg   env.Config
	store *db.PgxStore
	cache *redis.Client
}

func NewAPIServer(addr string, cfg env.Config, store *db.PgxStore, cache *redis.Client) *APIServer {
	return &APIServer{
		addr:  addr,
		cfg:   cfg,
		store: store,
		cache: cache,
	}
}

func (server *APIServer) Run() error {
	router := http.NewServeMux()

	// Sessions
	const authSessionTTL = time.Hour * 24 * 30
	const unauthSessionTTL = time.Hour
	sessionStore := cache.NewRedisCache(server.cache)
	sessionsHandler := sessions.NewSessionHandler(slog.Default(), sessionStore, authSessionTTL, unauthSessionTTL)

	// Authorization
	authProviders := make(map[string]auth.ProviderConfig)
	authProviders["google"] = auth.ProviderConfig{
		Issuer:       "https://accounts.google.com",
		ClientID:     server.cfg.OAUTH2_GOOGLE_CLIENT_ID,
		ClientSecret: server.cfg.OAUTH2_GOOGLE_CLIENT_SECRET,
		Scopes:       []string{"openid", "profile", "email"},
	}
	authHandler, err := auth.NewAuthHandler(slog.Default(), services.HandleHTTPError, sessionsHandler, server.store, server.cfg.BASE_URI, server.cfg.UI_URI, authProviders)
	if err != nil {
		return err
	}
	authHandler.RegisterRoutes(router)

	userHandler := user.NewUserHandler(slog.Default(), services.HandleHTTPError, sessionsHandler, server.store)
	userHandler.RegisterRoutes(router)

	businessHandler := business.NewBusinessHandler(slog.Default(), services.HandleHTTPError, sessionsHandler, server.store)
	businessHandler.RegisterRoutes(router)

	return http.ListenAndServe(server.addr, router)
}
