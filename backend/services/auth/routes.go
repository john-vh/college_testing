package auth

import (
	"context"
	"net/http"
	"net/url"

	"github.com/WilliamTrojniak/StudentTests/backend/models"
	"github.com/WilliamTrojniak/StudentTests/backend/services"
	"github.com/WilliamTrojniak/StudentTests/backend/util"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func (auth *AuthHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /auth/{provider}", auth.handleErr(auth.handleLogin))
	router.HandleFunc("GET /auth/{provider}/callback", auth.handleErr(auth.handleLoginCallback))
	router.HandleFunc("POST /auth/logout", auth.handleErr(auth.handleLogout))

	auth.logger.Info("Registered auth routes")
}

func (auth *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	provider := r.PathValue("provider")
	client, ok := auth.providers[provider]
	if !ok {
		return services.NewNotFoundServiceError(nil)
	}

	state, err := util.RandString(16)
	if err != nil {
		return err
	}
	nonce, err := util.RandString(16)
	if err != nil {
		return err
	}

	redirect := r.URL.Query().Get("redirect")
	setCallbackCookie(w, "state", state)
	setCallbackCookie(w, "nonce", nonce)
	setCallbackCookie(w, "redirect", redirect)
	http.Redirect(w, r, client.config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)

	return nil
}

func (auth *AuthHandler) handleLoginCallback(w http.ResponseWriter, r *http.Request) error {
	provider := r.PathValue("provider")
	client, ok := auth.providers[provider]
	if !ok {
		return services.NewNotFoundServiceError(nil)
	}

	state, err := r.Cookie("state")
	if err != nil || r.URL.Query().Get("state") != state.Value {
		auth.logger.Debug("State did not match")
		return services.NewInternalServiceError(err)
	}

	oauth2Token, err := client.config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		auth.logger.Debug("Failed to exchange for token")
		return services.NewInternalServiceError(err)
	}

	rawIdToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		auth.logger.Debug("Failed to retreive raw token")
		return services.NewInternalServiceError(nil)
	}

	verifier := client.provider.Verifier(&oidc.Config{ClientID: client.config.ClientID})

	idToken, err := verifier.Verify(context.TODO(), rawIdToken)
	if err != nil {
		auth.logger.Debug("Failed to verify raw token")
		return services.NewInternalServiceError(err)
	}

	nonce, err := r.Cookie("nonce")
	if err != nil || idToken.Nonce != nonce.Value {
		auth.logger.Debug("Nonce did not match")
		return services.NewInternalServiceError(err)
	}

	var claims models.OpenIDClaims

	userInfo, err := client.provider.UserInfo(context.TODO(), oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		auth.logger.Debug("Failed to retreive user info")
		return services.NewInternalServiceError(err)
	}

	if err := userInfo.Claims(&claims); err != nil {
		auth.logger.Debug("Failed to retreive user info claims")
		return services.NewInternalServiceError(err)
	}

	var userId *uuid.UUID
	session, err := auth.sessions.GetSession(r)
	if err == nil {
		userId = session.GetUserId()
	}

	if userId != nil {
		auth.logger.Debug("UserId is not nil", "user_id", userId)
		err := auth.LinkAccount(context.TODO(), userId, provider, &claims)
		if err != nil {
			auth.logger.Debug("Failed to link account")
			return err
		}
	} else {
		userId, err := auth.SaveAccount(context.TODO(), provider, &claims, &models.UserCreate{})
		if err != nil {
			auth.logger.Debug("Failed to create account")
			return err
		}
		_, err = auth.sessions.SetNewSession(w, r, userId)
		if err != nil {
			auth.logger.Debug("Failed to set session")
			return err
		}
	}

	redirect := ""
	if c, err := r.Cookie("redirect"); err == nil {
		redirect = c.Value
	}

	fullRedirect, err := url.JoinPath(auth.redirectBaseURI, redirect)
	if err != nil {
		return services.NewInternalServiceError(err)
	}

	http.Redirect(w, r, fullRedirect, http.StatusFound)

	return nil
}

func (auth *AuthHandler) handleLogout(w http.ResponseWriter, r *http.Request) error {
	_, err := auth.sessions.SetNewSession(w, r, nil)
	if err != nil {
		return err
	}

	return nil
}
