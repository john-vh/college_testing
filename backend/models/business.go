package models

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type BusinessStatus string

const (
	BUSINESS_STATUS_PENDING  BusinessStatus = "pending"
	BUSINESS_STATUS_ACTIVE   BusinessStatus = "active"
	BUSINESS_STATUS_DISABLED BusinessStatus = "disabled"
)

type businessMeta struct {
	Id        uuid.UUID      `json:"id" db:"id"`
	Status    BusinessStatus `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	LogoUrl   *string        `json:"logo_url" db:"logo_url"`
}

type BusinessUpdate struct {
	Name    string `json:"name" db:"name" validate:"required,min=3,max=64"`
	Desc    string `json:"desc" db:"description" validate:"required,min=8,max=256"`
	Website string `json:"website" db:"website" validate:"required,http_url"`
}

type BusinessCreate struct {
	BusinessUpdate
}

type BusinessOverview struct {
	businessMeta
	Name string `json:"name" db:"name"`
}

type Business struct {
	businessMeta
	BusinessCreate
	UserId uuid.UUID `json:"user_id" db:"user_id"`
}

type BusinessQueryParams struct {
	Status *BusinessStatus
	UserId *uuid.UUID
}

func (b *Business) URI(baseURL string) (string, error) {
	// TODO: Point to a specific business
	return url.JoinPath(baseURL, "account", "businesses")
}
