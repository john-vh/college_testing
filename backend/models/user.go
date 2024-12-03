package models

import (
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	USER_STATUS_ACTIVE   UserStatus = "active"
	USER_STATUS_BANNED   UserStatus = "banned"
	USER_STATUS_DISABLED UserStatus = "disabled"
)

type UserRole string

const (
	USER_ROLE_ADMIN UserRole = "admin"
	USER_ROLE_USER  UserRole = "user"
)

type acctInfo struct {
	Email         string `json:"email" db:"email"`
	Name          string `json:"name" db:"name"`
	EmailVerified bool   `json:"email_verified" db:"email_verified"`
}

type OpenIDClaims struct {
	acctInfo
	Id string `json:"sub" db:"id"`
}

type UserCreate struct {
}

type UserOverview struct {
	UserCreate
	Id                        uuid.UUID  `json:"id" db:"id"`
	NotifyApplicationUpdated  bool       `json:"notify_application_updated" db:"notify_application_updated"`
	NotifyApplicationReceived bool       `json:"notify_application_received" db:"notify_application_received"`
	CreatedAt                 time.Time  `json:"created_at" db:"created_at"`
	Status                    UserStatus `json:"status" db:"status"`
	acctInfo
}

type UserAccount struct {
	acctInfo
	Provider string `json:"provider" db:"provider"`
	// IsPrimary bool      `json:"is_primary" db:"is_primary"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	UserOverview
	Roles    []UserRole    `json:"roles" db:"roles" validate:"required,dive"`
	Accounts []UserAccount `json:"accounts" db:"accounts"`
}

func (u *User) HasRole(role UserRole) bool {
	return slices.Contains(u.Roles, role)
}

func (u *User) IsStudent() bool {
	// HACK: Need to improve student verification
	return slices.ContainsFunc(u.Accounts, func(ua UserAccount) bool {
		return strings.HasSuffix(ua.Email, ".edu")
	})
}
