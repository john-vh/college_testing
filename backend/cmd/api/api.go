package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/john-vh/college_testing/backend/cache"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/env"
	"github.com/john-vh/college_testing/backend/services"
	"github.com/john-vh/college_testing/backend/services/auth"
	"github.com/john-vh/college_testing/backend/services/business"
	"github.com/john-vh/college_testing/backend/services/notifications"
	"github.com/john-vh/college_testing/backend/services/sessions"
	"github.com/john-vh/college_testing/backend/services/user"
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

	backgroundServices := make([]services.BackgroundService, 0)

	// Mail
	mailClient := notifications.NewMailClient(server.cfg.MAIL_HOST, server.cfg.MAIL_PORT, server.cfg.MAIL_USER, server.cfg.MAIL_PASSWORD, slog.Default())
	notificationsService := notifications.NewNotificationService(mailClient, server.cfg.UI_URI, server.cfg.TEMPLATES_DIR, slog.Default())
	backgroundServices = append(backgroundServices, notificationsService)

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

	businessHandler := business.NewBusinessHandler(
		slog.Default(),
		services.HandleHTTPError,
		sessionsHandler,
		userHandler,
		notificationsService,
		server.store)
	businessHandler.RegisterRoutes(router)

	for _, service := range backgroundServices {
		service.Start()
	}

	res := http.ListenAndServe(server.addr, services.RequestLoggerMiddleWare(slog.Default())(services.CORSMiddleware(router)))

	for _, service := range backgroundServices {
		service.Stop()
	}
	return res
}
