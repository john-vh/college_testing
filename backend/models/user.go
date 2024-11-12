package models

import (
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
	Id        uuid.UUID  `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	Status    UserStatus `json:"status" db:"status"`
	acctInfo
}

type User struct {
	UserOverview
	Roles    []UserRole `json:"-" db:"roles" validate:"required,dive"`
	Accounts []struct {
		acctInfo
		Provider string `json:"provider" db:"provider"`
		// IsPrimary bool      `json:"is_primary" db:"is_primary"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	} `json:"accounts" db:"accounts"`
}
