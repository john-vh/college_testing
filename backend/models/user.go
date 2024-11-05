package models

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserStatus int

const (
	USER_STATUS_ACTIVE UserStatus = iota
	USER_STATUS_BANNED
	USER_STATUS_DISABLED
)

func (us UserStatus) String() string {
	switch us {
	case USER_STATUS_ACTIVE:
		return "active"
	case USER_STATUS_BANNED:
		return "banned"
	case USER_STATUS_DISABLED:
		return "disabled"
	default:
		return "unknown"
	}
}

func (role *UserStatus) ScanText(value pgtype.Text) error {
	switch value.String {
	case "active":
		*role = USER_STATUS_ACTIVE
		return nil
	case "banned":
		*role = USER_STATUS_BANNED
		return nil
	case "disabled":
		*role = USER_STATUS_DISABLED
		return nil
	default:
		return errors.New("Unsupported value scanning user status")
	}
}

func (role UserStatus) TextValue() (pgtype.Text, error) {
	val := pgtype.Text{}
	err := val.Scan(role.String())
	return val, err

}

type UserRole int

const (
	USER_ROLE_ADMIN UserRole = iota
)

func (ur UserRole) String() string {
	switch ur {
	case USER_ROLE_ADMIN:
		return "admin"
	default:
		return "unknown"
	}
}

func (role *UserRole) ScanText(value pgtype.Text) error {
	switch value.String {
	case "admin":
		*role = USER_ROLE_ADMIN
		return nil
	default:
		return errors.New("Unsupported value scanning user role")
	}
}

func (role UserRole) TextValue() (pgtype.Text, error) {
	val := pgtype.Text{}
	err := val.Scan(role.String())
	return val, err

}

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

type UserAccount struct {
	acctInfo
	Provider string `json:"provider" db:"provider"`
	// IsPrimary bool      `json:"is_primary" db:"is_primary"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	UserOverview
	Roles    []UserRole    `json:"-" db:"roles" validate:"required,dive"`
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
