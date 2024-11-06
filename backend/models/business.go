package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type BusinessStatus int

const (
	BUSINESS_STATUS_PENDING BusinessStatus = iota
	BUSINESS_STATUS_ACTIVE
	BUSINESS_STATUS_DISABLED
)

func (bs BusinessStatus) String() string {
	switch bs {
	case BUSINESS_STATUS_PENDING:
		return "pending"
	case BUSINESS_STATUS_ACTIVE:
		return "active"
	case BUSINESS_STATUS_DISABLED:
		return "disabled"
	default:
		return "unknown"
	}
}

func (s *BusinessStatus) ScanText(value pgtype.Text) error {
	switch value.String {
	case "pending":
		*s = BUSINESS_STATUS_PENDING
		return nil
	case "active":
		*s = BUSINESS_STATUS_ACTIVE
		return nil
	case "disabled":
		*s = BUSINESS_STATUS_DISABLED
		return nil
	default:
		return errors.New("Unsupported value scanning business status")
	}
}

func (s BusinessStatus) TextValue() (pgtype.Text, error) {
	val := pgtype.Text{}
	err := val.Scan(s.String())
	return val, err
}

type BusinessUpdate struct {
	Website string `json:"website" db:"website" validate:"required,http_url"`
	Desc    string `json:"desc" db:"description" validate:"required,min=8,max=256"`
}

type BusinessCreate struct {
	Name string `json:"name" db:"name" validate:"required,min=3,max=64"`
	BusinessUpdate
}

type Business struct {
	BusinessCreate
	Id        uuid.UUID      `json:"id" db:"id"`
	UserId    uuid.UUID      `json:"-" db:"user_id"`
	Status    BusinessStatus `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}

type BusinessQueryParams struct {
	Status *BusinessStatus
	UserId *uuid.UUID
}
