package sessions

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/WilliamTrojniak/StudentTests/backend/cache"
	"github.com/WilliamTrojniak/StudentTests/backend/services"
	"github.com/WilliamTrojniak/StudentTests/backend/util"
	"github.com/google/uuid"
)

const (
	session_cookie_name = "session"
	csrf_header_name    = "X-CSRF-TOKEN"
)

type sessionData struct {
	UserId    *uuid.UUID
	CSRFToken string
}

type Session struct {
	Id   string
	ttl  time.Duration
	data sessionData
}

func (s *Session) GetUserId() *uuid.UUID {
	return s.data.UserId
}

type SessionsHandler struct {
	logger          *slog.Logger
	store           cache.Cache
	authorizedTTL   time.Duration
	unauthorizedTTL time.Duration
}

func NewSessionHandler(logger *slog.Logger, store cache.Cache, authedTTL time.Duration, unauthedTTL time.Duration) *SessionsHandler {
	if logger == nil {
		logger = slog.Default()
	}

	return &SessionsHandler{
		logger:          logger,
		store:           store,
		authorizedTTL:   authedTTL,
		unauthorizedTTL: unauthedTTL,
	}
}

func (h *SessionsHandler) SetNewSession(w http.ResponseWriter, r *http.Request, userId *uuid.UUID) (*Session, error) {
	newSession, err := h.newSessionFromUserId(userId)
	if err != nil {
		return nil, err
	}

	err = h.saveSessionToStore(context.TODO(), newSession)
	if err != nil {
		return nil, err
	}

	oldSessionId, err := h.getSessionIdFromCookie(r)
	if err == nil { // Check if an old session exists
		err = h.store.Delete(context.TODO(), oldSessionId)
		if err != nil {
			h.logger.Warn("Failed to delete old session.", "sessionId", oldSessionId)
		}
	}

	h.saveSessionToResponse(w, newSession)

	return newSession, nil
}

func (h *SessionsHandler) GetSession(r *http.Request) (*Session, error) {
	sessionId, err := h.getSessionIdFromCookie(r)
	if err != nil {
		return nil, err
	}

	session, err := h.getSessionFromStore(r.Context(), sessionId)
	if err != nil {
		return nil, services.NewUnauthenticatedServiceError(err)
	}

	return session, nil
}

func (h *SessionsHandler) newSessionFromUserId(userId *uuid.UUID) (*Session, error) {
	csrftoken, err := util.RandString(16)
	if err != nil {
		return nil, err
	}

	sessionId, err := util.RandString(32)
	if err != nil {
		return nil, err
	}

	ttl := h.authorizedTTL
	if userId == nil {
		ttl = h.unauthorizedTTL
	}

	return &Session{
		Id:  sessionId,
		ttl: ttl,
		data: sessionData{
			UserId:    userId,
			CSRFToken: csrftoken,
		},
	}, nil
}

func (h *SessionsHandler) saveSessionToStore(ctx context.Context, session *Session) error {
	data, err := json.Marshal(session.data)
	if err != nil {
		return err
	}

	err = h.store.Set(ctx, session.Id, data, session.ttl)
	if err != nil {
		return err
	}

	return nil
}

func (h *SessionsHandler) saveSessionToResponse(w http.ResponseWriter, session *Session) {
	c := &http.Cookie{
		Name:     session_cookie_name,
		Value:    session.Id,
		MaxAge:   int(session.ttl.Seconds()),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: 4,
	}
	w.Header().Set(csrf_header_name, session.data.CSRFToken)
	http.SetCookie(w, c)
}

func (h *SessionsHandler) getSessionIdFromCookie(r *http.Request) (string, error) {
	c, err := r.Cookie(session_cookie_name)
	if err != nil {
		return "", services.NewUnauthenticatedServiceError(err)
	}
	return c.Value, nil
}

func (h *SessionsHandler) getSessionFromStore(ctx context.Context, id string) (*Session, error) {

	data, err := h.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	session := &Session{
		Id:  id,
		ttl: h.authorizedTTL,
	}

	err = json.Unmarshal(data, &session.data)
	if err != nil {
		return nil, err
	}

	if (*session.data.UserId) == uuid.Nil {
		session.data.UserId = nil
		session.ttl = h.unauthorizedTTL
	}

	return session, nil
}

func (h *SessionsHandler) deleteSessionFromStore(ctx context.Context, session *Session) error {
	return h.store.Delete(ctx, session.Id)
}
